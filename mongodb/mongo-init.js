db = db.getSiblingDB("schedule_db");

db.createCollection("users");
db.createCollection("classrooms");
db.createCollection("courses");
db.createCollection("draft_schedules");
db.createCollection("previous_schedules");
