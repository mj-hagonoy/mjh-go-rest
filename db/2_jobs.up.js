//use codev;
db = db.getSiblingDB('codev')
db.createCollection( "jobs", {
    validator: { $jsonSchema: {
        bsonType: "object",
        required: ["type", "status", "initiated_by"],
        properties: {
            type : { bsonType: "string"},
            status: { bsonType: "string"},
            initiated_by: {bsonType: "string"},
            source_file: {bsonType: "string"},
            created: {bsonType: "string"},
            modified: {bsonType: "string"},
        }
    }}
})
