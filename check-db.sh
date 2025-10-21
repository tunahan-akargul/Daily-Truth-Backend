#!/bin/bash
# Simple script to check MongoDB data

echo "=========================================="
echo "Checking MongoDB Collections"
echo "=========================================="
echo ""

echo "Database: main-services"
echo ""

echo "Collections:"
mongosh main-services --quiet --eval 'db.getCollectionNames()'
echo ""

echo "Words count:"
mongosh main-services --quiet --eval 'db.words.countDocuments({})'
echo ""

echo "All words:"
mongosh main-services --quiet --eval 'db.words.find().forEach(printjson)'

