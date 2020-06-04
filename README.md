# MySQL DB Manager
web interface for managing local a database
Note: **DO NOT USE THIS IN PRODUCTION** There is no safety or security build in.



## TODO:
 - cleanup code
 - make form async
 - show query time
 - save query
  - make adjustable params (form)
  - decide how to show/run
 - join queries into tree, somehow (orders -> order_items)
  - not sure how to do this in go, need query keys, maybe I can use "describe"
 - fix "Describe" button
  - don't show if query starts with "describe"
 - show all
  - override limit
 - multiple sessions (use defined color/theme)
  - a. store db connection in session?
  - b. add db connection id to path e.g. `/root@127.0.0.1:3306/db01/run?query=...`
 - ssh connect
  - SSH Hostname: 25.25.25.25:22
  - SSH Username: root
  - SSH Password: 
  - SSH Key File: /Users/jaredwarren/.ssh/id_rsa
  - MySQL Hostname: 25.25.25.25
  - MySQL Server Port: 3306
  - Username: root
  - password: 12345
  - Default Schema: db01
 - make native app? (more ui space) (multiple windows?)
 - update ui to match current query
  - describe
  - sort
  - show_all
  - ...
  - is there a "generic" way to do this e.g. `$('#describe')->hide()` or `$('#sort_by_user_id').disabled()`
 - test connection on login?