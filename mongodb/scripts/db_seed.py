import pandas as pd  
from pymongo import MongoClient 
from dotenv import load_dotenv
from pathlib import Path 
import os


def db_seed(): 

    classrooms_df = pd.read_csv("../data/classrooms.csv") 
    courses_df = pd.read_csv("../data/courses.csv") 
    users_df = pd.read_csv("../data/users.csv")  

    dotenv_path = Path('../../.env')
    load_dotenv(dotenv_path=dotenv_path)  

    mongo_user = os.getenv("MONGO_USERNAME")
    mongo_pw = os.getenv("MONGO_PASSWORD")
    mongo_ip = os.getenv("MONGO_ADDRESS")
    mongo_port = os.getenv("MONGO_PORT") 

    conn = MongoClient('mongodb://'+mongo_user+':'+mongo_pw+'@'+mongo_ip+':'+mongo_port+'/')
    db = conn['schedule_db']


    print(classrooms_df.head()) 
    print("################### \n") 
    print(courses_df.head())  
    print("################### \n")   
    print(users_df.head()) 

if __name__ == "__main__": 
    db_seed()
