//use codev;
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

//db.users.createIndex({"_id": 1})
db.users.createIndex({"type": 1}, {unique:true, sparse:true})
db.users.createIndex({"type": 1, "initiated_by": 1}, {unique:true, sparse:true})