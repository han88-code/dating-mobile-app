# dating-mobile-app
In this Dating Mobile App project provide some basic backend system functionalities that similar to Tinder or Bumble.
- User able to register, login & logout.
- User able to view, swipe left (pass) & swipe right (like) 10 other user (dating) profiles in total (pass + like) in 1 day.
- Same profiles can’t appear twice in the same day.
- User able to set to verified status, so that the user have no swipe quota (can see more than 10 other user profiles).


--- 
<h2>Requirement</h2>

1. Go
2. Mysql
3. Postman

--- 
<h2>Structure Folder</h2>

<pre>
-> dating-mobile-app
  -> .env
  -> main.go
  -> go.mod
  -> go.sum
  -> app
      -> config
          -> config.go
      -> controllers
          -> user.go
      -> models
          -> user.go
      -> routes
          -> routes.go
      -> utils
          -> utils.go
</pre>

--- 
<h2>Preparation</h2>

1. Database<br />
   a. Create database name: [_dating_mobile_db_].<br />
   b. Import file [_user.sql_] & [_user_logs.sql_] from [**_/resource/database/_**] folder.
2. Postman<br />
   a. Import file [_Local.postman_environment.json_] & [_User.postman_collection.json_] from [**_/resource/postman/_**] folder.

--- 
<h2>Setting</h2>

1. Rename [.env.example] file to [.env]
2. Update the variable value based on your local configuration.
3. Open terminal (or command prompt) on your project directory path, then run this command line <pre>go run main.go</pre>

--- 
<h2>API Collection</h2>

1. Register API - For register new user.
2. Login API - For login existing user.
3. Home API - For get login user information and other user profile information.
4. Swipe API - For get next other user profile information.
5. Logout API - For logout existing user.
