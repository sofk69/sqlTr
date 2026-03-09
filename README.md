A web application for managing users and their payments, written in Go. 
It allows you to create users, add payments, and search for users by name.

Features:
1) User Management: create new users with automatic data generation
2) Payment Management: add payments for existing users
3) Search:find users by first name
4) Data Generation: automatic generation of random names, emails, ages, and payment names
5) Web Interface: user-friendly interface for interacting with the system

Technologies:
Backend: Go 1.25+
Database: PostgreSQL
Database Driver: lib/pq
Templating: HTML templates
Frontend: CSS

Data Models
User:
ID - unique identifier
Age - age
FirstName - first name
LastName - last name
Email - email address

Payment:
PaymentID - unique payment identifier
UserID - user ID
Amount - payment amount
PaymentName - payment name
Installation and Setup

Prerequisites:
Go 1.25 or higher
PostgreSQL
Git
