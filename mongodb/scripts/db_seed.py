import pandas as pd 


def db_seed(): 

    classrooms_df = pd.read_csv("../data/classrooms.csv") 
    courses_df = pd.read_csv("../data/courses.csv") 
    users_df = pd.read_csv("../data/users.csv") 

    print(classrooms_df.head()) 
    print("################### \n") 
    print(courses_df.head())  
    print("################### \n")   
    print(users_df.head()) 

if __name__ == "__main__": 
    db_seed()
