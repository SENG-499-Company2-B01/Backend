import pandas as pd  
from pymongo import MongoClient 
from dotenv import load_dotenv
import os  
import bcrypt 


def check_if_not_empty(coll,coll_name):  

    if coll.count_documents({}) > 0: 

        print("INFO: Data found in %s collection ... skipping inserting data in this collection" %(coll_name)) 
        return True  
    
    print("INFO: No data found in %s collection ... inserting data" %(coll_name))
    return False 

def gen_pw(input,secret): 

    input = input + str(secret)
    bytes = input.encode('utf-8') 
    salt = bcrypt.gensalt() 
    return bcrypt.hashpw(bytes,salt) 

def parse_prereq(input): 

    input = input.replace('[','').replace(']','') 

    if input == '': 
        return [['']] 
    
    output = input.split(",") 
    for i in range(len(output)): 
        output[i] = output[i].split(":") 

    return output 

def add_course(schedule,row,term_count,df,index): 

    course = row['Course']  
    section_num = row['Num'] 
    building = row['Building'] 
    professor = row['Professor'] 
    num_seats = int(row['Num_Seats'])
    num_registered = int(row['Enrolled'])

    # We simply do not have these values currently in our dataset
    start_time = "" 
    end_time = "" 

    course_dict = {} 
    course_dict['course'] = course 
    course_dict['sections'] = [{"num":section_num,"building":building, 
                                "professor":professor,"days":[], 
                                "num_seats":num_seats,"num_registered":num_registered, 
                                "start_time":start_time,"end_time":end_time}] 
    
    # Check if there is an A02 section, and if so add it to the sections JSON
    try:
        new_row = df.iloc[index+1]  

        if str(new_row['Num']) == "A02": 

            course = new_row['Course']  
            section_num = new_row['Num'] 
            building = new_row['Building'] 
            professor = new_row['Professor'] 
            num_seats = int(new_row['Num_Seats'])
            num_registered = int(new_row['Enrolled'])

            course_dict['sections'].append({"num":section_num,"building":building, 
                                "professor":professor,"days":[], 
                                "num_seats":num_seats,"num_registered":num_registered, 
                                "start_time":start_time,"end_time":end_time})


    except:  
        print("INFO ... End of classes.csv file")


    schedule["terms"][term_count]["courses"].append(course_dict) 

    return schedule


def load_users(coll): 

    if check_if_not_empty(coll,'users'): 
        return   
    
    admin_1 = os.getenv("ADMIN_1") 
    admin_2 = os.getenv("ADMIN_2") 
    pw_secret = os.getenv("PW_SECRET")


    users_df = pd.read_csv("users.csv").fillna('')
    users_df = users_df.astype({'Credentials':'string'}) 
    pref_dict = {"M":[["08:30","16:00"]],"T":[["08:30","16:00"]],"W":[["08:30","16:00"]],"R":[["08:30","16:00"]],"F":[["08:30","16:00"]]}

    for index, row in users_df.iterrows():

        user = {} 
        user['username'] = row['Firstname'] + '.' + row['Lastname'] 
        user['email'] = row['Email'] 
        user['password'] = gen_pw(user['username'],pw_secret) 
        user['name'] = row['Firstname'] + ' ' + row['Lastname']

        if user['username'].lower() == admin_1 or user['username'].lower() == admin_2:
            user['isAdmin'] = True
        else: 
            user['isAdmin'] = False 

        user['peng'] = (row['Peng'] == 1) 
        user['pref_approved'] = False
        user['max_courses'] = 6
        # Convert qualifications string to array of strings 
      
        string_qualifications = row['Credentials'].replace("[","").replace("]","")
        qualifications = string_qualifications.split(",")
        user['course_pref'] = qualifications 
        user['time_pref'] = pref_dict  
        user['available'] = pref_dict

        coll.insert_one(user)

    return


def load_courses(coll):  

    if check_if_not_empty(coll,'courses'):
        return  

    courses_df = pd.read_csv("courses.csv")  

    for index, row in courses_df.iterrows(): 
        
        course = {} 
        course['shorthand'] = row['Course'] 
        course['name'] = row['Name'] 
        course['terms_offered'] = row['Offered'].replace("[","").replace("]","").split(",")
        course['prerequisites'] = parse_prereq(row['Prerequisites']) 

        coll.insert_one(course) 

    return 


def load_classrooms(coll): 

    if check_if_not_empty(coll,'classrooms'): 
        return   

    classrooms_df = pd.read_csv("classrooms.csv")  

    for index, row in classrooms_df.iterrows(): 

        classroom = {}
        classroom['building'] = row['Shorthand'] 
        classroom['capacity'] = row['Capacity'] 
        classroom['room'] = row['Room Number'] 

        coll.insert_one(classroom)

    return 

def load_old_schedules(coll): 

    if check_if_not_empty(coll,'prev_schedules'): 
        return 
    
    schedules_df = pd.read_csv("classes.csv") 
    schedule = {} 
    prev_year = 2008 
    prev_term = "summer" 
    schedule['year'] = prev_year 
    schedule['terms'] = [{"term":prev_term,"courses":[]}] 
    term_count = 0

    for index, row in schedules_df.iterrows(): 

        year = row['Year']  
        term = row['Term']
        
        if prev_year != year: 
            coll.insert_one(schedule) 
            schedule = {} 
            schedule["year"] = year  
            schedule["terms"] = [{"term":term,"courses":[]}]
            prev_year = year  
            prev_term = term
            term_count = 0

        
        if prev_term != term:  
            term_count += 1 
            schedule["terms"].append({"term":term,"courses":[]})   

        # Skip A02 Terms since we already have added them via the add_courses function
        if str(row['Num']) == "A02": 
            continue


        schedule = add_course(schedule,row,term_count,schedules_df,index) 
    
    coll.insert_one(schedule)

def db_seed():

    load_dotenv()

    mongo_user = os.getenv("MONGO_LOCAL_USERNAME")
    mongo_pw = os.getenv("MONGO_LOCAL_PASSWORD")
    mongo_host = os.getenv("MONGO_LOCAL_HOST") 

    conn = MongoClient('mongodb://'+mongo_user+':'+mongo_pw+'@'+mongo_host+'/')
    db = conn['schedule_db']

    user_collection = db['users']  
    courses_collection = db['courses'] 
    classrooms_collection = db['classrooms'] 
    schedules_collection = db['previous_schedules']


    load_users(user_collection)
    load_courses(courses_collection) 
    load_classrooms(classrooms_collection) 
    load_old_schedules(schedules_collection)

if __name__ == "__main__": 
    db_seed()
