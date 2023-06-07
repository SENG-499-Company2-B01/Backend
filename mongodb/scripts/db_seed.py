import pandas as pd  
from pymongo import MongoClient 
from dotenv import load_dotenv
from pathlib import Path 
import os 

def check_if_not_empty(coll,coll_name):  

    if coll.count_documents({}) > 0: 

        print("INFO: Data found in %s collection ... skipping inserting data in this collection" %(coll_name)) 
        return True  
    
    print("INFO: No data found in %s collection ... inserting data" %(coll_name))
    return False

def load_users(coll,admin_1,admin_2): 

    if check_if_not_empty(coll,'users'): 
        return 


    users_df = pd.read_csv("../data/users.csv") 

    for index, row in users_df.iterrows():

        user = {} 
        user['username'] = row['Firstname'] + '.' + row['Lastname'] 
        user['email'] = row['Email'] 
        user['password'] = '' 
        user['firstname'] = row['Firstname'] 
        user['lastname'] = row['Lastname'] 

        if user['username'].lower() == admin_1 or user['username'].lower() == admin_2:
            user['isAdmin'] = True
        else: 
            user['isAdmin'] = False

        user['prefrences'] = []
        user['qualifications'] = row['Credentials']  

        coll.insert_one(user)

    return


def load_courses(coll):  

    if check_if_not_empty(coll,'courses'):
        return  

    courses_df = pd.read_csv("../data/courses.csv")  

    for index, row in courses_df.iterrows(): 

        course = {} 
        course['shorthand'] = row['Course'] 
        course['name'] = row['Name'] 
        course['offered'] = row['Offered'] 
        course['equipment'] = [] 
        course['prerequisites'] = [] 

        coll.insert_one(course) 

    return 


def load_classrooms(coll): 
    if check_if_not_empty(coll,'classrooms'): 
        return   

    classrooms_df = pd.read_csv("../data/classrooms.csv")  

    for index, row in classrooms_df.iterrows(): 

        classroom = {} 
        classroom['shorthand'] = row['Shorthand'] 
        classroom['building'] = row['Building Name'] 
        classroom['capacity'] = row['Capacity'] 
        classroom['room_number'] = row['Room Number'] 
        classroom['Equipment'] = []  

        coll.insert_one(classroom)

    return


def db_seed():

    dotenv_path = Path('../../.env')
    load_dotenv(dotenv_path=dotenv_path)  

    mongo_user = os.getenv("MONGO_USERNAME")
    mongo_pw = os.getenv("MONGO_PASSWORD")
    mongo_ip = os.getenv("MONGO_ADDRESS")
    mongo_port = os.getenv("MONGO_PORT")  
    admin_1 = os.getenv("ADMIN_1") 
    admin_2 = os.getenv("ADMIN_2")

    conn = MongoClient('mongodb://'+mongo_user+':'+mongo_pw+'@'+mongo_ip+':'+mongo_port+'/')
    db = conn['schedule_db']

    user_collection = db['users']  
    courses_collection = db['courses'] 
    classrooms_collection = db['classrooms']


    load_users(user_collection,admin_1,admin_2) 
    load_courses(courses_collection) 
    load_classrooms(classrooms_collection)

if __name__ == "__main__": 
    db_seed()
