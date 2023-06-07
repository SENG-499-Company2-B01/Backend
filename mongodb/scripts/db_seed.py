import pandas as pd  
from pymongo import MongoClient 
from dotenv import load_dotenv
from pathlib import Path 
import os


ADMIN_1 = 'rich.little' 
ADMIN_2 = 'dan.mai'

def load_users(coll):  

    users_df = pd.read_csv("../data/users.csv") 

    for index, row in users_df.iterrows():

        user = {} 
        user['username'] = row['Firstname'] + '.' + row['Lastname'] 
        user['email'] = row['Email'] 
        user['password'] = '' 
        user['firstname'] = row['Firstname'] 
        user['lastname'] = row['Lastname'] 

        if user['username'].lower() == ADMIN_1 or user['username'].lower() == ADMIN_2:
            user['isAdmin'] = True
        else: 
            user['isAdmin'] = False

        user['prefrences'] = []
        user['qualifications'] = row['Credentials']  

        coll.insert_one(user)

    return 


# type Course struct {
# 	ShortHand     string   `json:"shorthand" bson:"shorthand"`
# 	Name          string   `json:"name" bson:"name"`
# 	Equipements   []string `json:"equipment" bson:"equipements"`
# 	Prerequisites []string `json:"prerequisites" bson:"prerequisites"`
# }

def load_courses(coll):  

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


def load_classrooms():  

    classrooms_df = pd.read_csv("../data/classrooms.csv") 


    return


def db_seed():

    dotenv_path = Path('../../.env')
    load_dotenv(dotenv_path=dotenv_path)  

    mongo_user = os.getenv("MONGO_USERNAME")
    mongo_pw = os.getenv("MONGO_PASSWORD")
    mongo_ip = os.getenv("MONGO_ADDRESS")
    mongo_port = os.getenv("MONGO_PORT") 

    conn = MongoClient('mongodb://'+mongo_user+':'+mongo_pw+'@'+mongo_ip+':'+mongo_port+'/')
    db = conn['schedule_db']  

    user_collection = db['users']  
    courses_collection = db['courses']

    load_users(user_collection) 
    load_courses(courses_collection)

if __name__ == "__main__": 
    db_seed()
