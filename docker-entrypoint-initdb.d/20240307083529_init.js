// source データベース
db = db.getSiblingDB("source")

db.users.drop()
db.createCollection("users", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["_id", "name"],
    },
  },
  validationLevel: "moderate",
})
const aliceId = new ObjectId()
const bobId = new ObjectId()
const charlieId = new ObjectId()
db.users.insertMany([
  { _id: aliceId, name: "Alice" },
  { _id: bobId, name: "Bob" },
  { _id: charlieId, name: "Charlie" },
])

db.points.drop()
db.createCollection("points", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["_id", "userId", "point"],
      properties: {
        point: {
          bsonType: ["int", "long"],
          minimum: 0,
        },
      },
    },
  },
  validationLevel: "moderate",
})
db.points.insertMany([
  { _id: new ObjectId(), userId: aliceId, point: 100 },
  { _id: new ObjectId(), userId: aliceId, point: 200 },
  { _id: new ObjectId(), userId: aliceId, point: 300 },
  { _id: new ObjectId(), userId: bobId, point: 200 },
  { _id: new ObjectId(), userId: bobId, point: 400 },
  { _id: new ObjectId(), userId: bobId, point: 600 },
  { _id: new ObjectId(), userId: charlieId, point: 300 },
  { _id: new ObjectId(), userId: charlieId, point: 500 },
  { _id: new ObjectId(), userId: charlieId, point: 700 },
])

// aggregate データベース
db = db.getSiblingDB("aggregate")
db.leaderboard.drop()
db.createCollection("leaderboard", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["_id", "userId"],
    },
  },
  validationLevel: "moderate",
})
db.leaderboard.createIndex({ userId: 1, method: 1 }, { unique: true })
