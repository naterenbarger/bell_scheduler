Title: Scheduler Application - Web Interface with Go Backend and Vue.js Frontend

1. Overview: This document outlines the requirements for a web-based application that allows users to manage schedules based on specific days of the week, times when relays should trigger in HH:MM format, a global configurable duration for triggering these relay rings, and an option to set one schedule as default.

2. Scope:

Frontend: Develop a user-friendly interface for creating, reading, updating, deleting schedules.
Backend: Provide API endpoints for managing schedules and serve static files from the frontend application.
User Authentication and Management Interface: Implement secure user authentication with an initial administrator account to manage users.
3. Technical Requirements:

3.1 Frontend Framework:
The Vue.js framework will be used to build the user interface of the application.

Version Requirement:

Vue.js version 2.x or later
Build Tooling:

Node.js (v14.x or later)
npm or yarn for package management
Build tools like Webpack, Vite, etc., as per your preference
3.2 Backend Framework:
The Go language and Gin framework will be used to build the backend.

Go Version:

Go version 1.17.x or later
Gin Framework:

Dependency management using go get for Gin.
Authentication Libraries:

Use libraries like JWT (JSON Web Tokens) for secure user authentication and authorization.
3.3 Frontend Static Files:
The Vue.js application will be built and served statically from the backend server using Go's file handling capabilities.

4. Functionality Requirements:

User Authentication with Initial Admin Account:

An initial administrator account with a default username (e.g., admin) and password (to be specified later).
User login to access protected endpoints.
Logout functionality for user sessions.
CRUD Interface for Managing Users:

List all users, including their roles.
Create new users with username, email, password, and role.
Update existing user details (username, email, role).
Delete a user (only by an admin).
Global Configurable Bell Ring Duration:

Allow users to set a global configurable duration for how long each bell ringing trigger should last. This setting will be applied to all scheduled triggers unless specified otherwise at a later date.
Default Schedule Option:

Each schedule can have an option to mark it as the default.
When marked, this default schedule will be automatically activated at midnight every day.
Schedule Structure Breakdown:

A schedule consists of two parts:
Schedule Information: Contains the following details:
Schedule name
Description (optional)
Default flag indicating whether it is the active schedule
Schedule Times: Includes multiple time slots, each with the following details:
Time slot ID
Trigger time in HH:MM format
Days of the week (e.g., Monday, Tuesday)
Description (optional)
CRUD Interface for Schedules:

Create a new schedule with the following details:
Schedule name
Description (optional)
Default flag indicating whether it should be activated as default
Read existing schedules.
Update an existing schedule to modify the schedule name, description, and default status.
Delete an existing schedule.
CRUD Interface for Schedule Times:

Add time slots to an existing schedule with the following details:
Trigger time in HH:MM format
Days of the week (e.g., Monday, Tuesday)
Description (optional)
Update and delete existing time slots as needed.
5. User Interface Requirements:

Login Page:

Allow users to enter their username and password.
Display a login form with fields for username and password.
Include a submit button for user authentication.
User Dashboard:

Upon successful login, display the user dashboard where they can manage schedules, view global settings, and switch between active schedules if needed.
Provide links to navigate to different sections such as schedule management, settings, and account management.
Schedule Management Page:

Display a list of all existing schedules with options to create new ones.
Include search functionality to filter schedules based on name or description.
Allow users to edit an existing schedule by clicking on the schedule entry.
Provide a confirmation dialog for deleting an existing schedule.
Schedule Detail Page:

Display details of a selected schedule including its default status and list of time slots.
Include options to add new time slots or edit/delete existing ones.
Allow users to toggle the default flag for a schedule from this page.
Settings Page:

Provide global settings such as configuring the bell ring duration that affects all schedules unless overridden in individual schedules.
Display current settings and allow updates through input fields.
6. Scheduling Logic Requirements:

Default Schedule Activation:

At midnight every day, check if any schedule is marked as default and activate it.
Ensure that the active schedule remains unchanged until another default schedule is activated or a user manually changes it.
Time Slot Validation:

Validate time slots to ensure they do not overlap within the same day.
Enforce constraints such as valid HH:MM formats for trigger times.
7. Security Requirements:

Authentication and Authorization:

Implement JWT-based authentication with roles (e.g., admin, user) to restrict access based on user privileges.
Ensure secure transmission of credentials using HTTPS.
Data Validation:

Validate all input data from users to prevent injection attacks and ensure data integrity.
8. Error Handling Requirements:

User-Friendly Errors:

Display clear, user-friendly error messages for common issues such as authentication failures or invalid inputs.
Provide guidance on how to resolve these errors.
Logging:

Implement logging for all critical operations and errors for debugging and monitoring purposes.
9. Performance Requirements:

Responsive Design:
Ensure the application is responsive across different devices (desktop, tablet, mobile).
Optimizations:
Optimize database queries to improve performance.
Use caching where appropriate to reduce latency.
This breakdown ensures that the application is well-structured, secure, scalable, and provides a user-friendly interface with clear navigation options, including the ability to set a default schedule and have it applied automatically at midnight each day.