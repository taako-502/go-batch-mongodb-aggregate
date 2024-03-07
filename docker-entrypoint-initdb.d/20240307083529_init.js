db = db.getSiblingDB("source")

db.users.drop()
const aliceId = new ObjectId()
const bobId = new ObjectId()
const charlieId = new ObjectId()
db.users.insertMany([
  { _id: aliceId, name: "Alice" },
  { _id: bobId, name: "Bob" },
  { _id: charlieId, name: "Charlie" },
])

db.points.drop()
db.points.insertMany([
  { _id: new ObjectId(), userId: aliceId, point: 100 },
  { _id: new ObjectId(), userId: aliceId, point: 200 },
  { _id: new ObjectId(), userId: aliceId, point: 300 },
  { _id: new ObjectId(), userId: bobId, point: 200 },
  { _id: new ObjectId(), userId: bobId, point: 400 },
  { _id: new ObjectId(), userId: charlieId, point: 300 },
  { _id: new ObjectId(), userId: charlieId, point: 500 },
  { _id: new ObjectId(), userId: charlieId, point: 600 },
])

db = db.getSiblingDB("aggregate")
db.leaderboard.drop()
db.leaderboard.createIndex({ userId: 1 }, { unique: true })
