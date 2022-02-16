//use codev;'
db = db.getSiblingDB('codev')
db.createCollection( "users", {
    validator: { $jsonSchema: {
        bsonType: "object",
        required: ["first_name", "last_name", "email"],
        properties: {
            first_name : { bsonType: "string"},
            last_name: { bsonType: "string"},
            email: {bsonType: "string"},
            password: {bsonType: "string"},
            activated: {bsonType: "bool"},
            created: {bsonType: "string"},
            modified: {bsonType: "string"},
        }
    }}
})

//db.users.createIndex({"_id": 1})
db.users.createIndex({"email": 1}, {unique:true, sparse:true})
db.users.createIndex({"first_name": 1, "last_name": 1}, {unique:true, sparse:true})