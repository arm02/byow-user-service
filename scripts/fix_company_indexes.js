// MongoDB script to fix company email index with sparse option
// Run this in MongoDB shell or MongoDB Compass

// Connect to your database
use('byow-user-service');

// Drop the existing problematic index
db.companies_collections.dropIndex("company_email_unique");

// Create new sparse unique index for company email
db.companies_collections.createIndex(
  { "email": 1 },
  { 
    "unique": true, 
    "sparse": true, 
    "name": "company_email_unique" 
  }
);

// Verify the index was created
db.companies_collections.getIndexes();

print("Company email index fixed successfully!");