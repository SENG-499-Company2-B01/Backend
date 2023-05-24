db = db.getSiblingDB("admin") 
db.auth("admin","admin") 
db = db.getSiblingDB("schedule_db") 

db.createUser({ 
    user: "user", 
    pwd: "user", 
    roles: [ 
        { 
            role: "readWrite", 
            db: "schedule_db"
        }
    ]
}); 

db.createCollection("users") 
db.createCollection("classrooms") 
db.createCollection("courses") 
db.createCollection("schedules")