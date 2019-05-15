db.orders.drop();
// db.container.createIndex({ myfield: 1 }, { unique: true }),
db.orders.insertMany([
  {
    _id: "5cda21b1a94b7acd5d522973",
    id: "f0a31480-a03d-47d4-a7ac-e6c0a2d70123",
    project: "project name",
    items: [
      "consequat",
      "incididunt",
      "consequat",
      "amet",
      "ea",
      "eu",
      "Lorem",
      "consectetur",
      "mollit",
      "ex",
      "incididunt",
      "minim"
    ]
  },
  {
    _id: "5cda21c7b8dee979a4161227",
    id: "6ea5c742-8d47-47f4-936d-1535b1a08676",
    project: "project name",
    items: [
      "incididunt",
      "officia",
      "magna",
      "ut",
      "fugiat",
      "ipsum",
      "exercitation",
      "dolore",
      "ex"
    ]
  },
  {
    _id: "5cda21cf1fb0ecc0892b9075",
    id: "72c0eb46-3229-4536-b803-dd327160b5fe",
    project: "project name",
    items: [
      "eu",
      "exercitation",
      "irure",
      "reprehenderit",
      "aute",
      "ut",
      "adipisicing",
      "consectetur",
      "aliquip",
      "aliquip",
      "nulla",
      "nisi",
      "id"
    ]
  },
  {
    _id: "5cda23bbcde6bb894cc7f473",
    id: "aeb0b4cc-487f-4ad7-aab7-34be5b683599",
    project: "project name",
    items: ["deserunt", "duis", "sit", "in"]
  }
]);
