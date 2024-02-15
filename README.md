# Containers
This is a quick guide on how to get the backend and frontend up and running

## Download and run JSON server
1. Download and install [JSON Server](https://github.com/typicode/json-server)
   ```bash
   npm i -g json-server
   ```

2. Run
   ```bash
   json-server database.json
   ```

## Start the backend server
1. Download and install [go](https://go.dev/dl/)
2. Navigate into the `backend` folder
3. Run the server:[^1]
   ```bash
   # This will automatically compile and run the server
   go run backend
   ```
   
## Start the frontend server
1. Download and install [AngularCLI](https://angular.io)
   ```bash
   npm i -g @angular/cli
   ```
2. Navigate into the `frontend` folder
3. Run the server:[^2]
   ```bash
   ng serve
   # It will tell you the url in the terminal after it started
   ```

---

[^1]: Or just run the executable in the backend directory
[^2]: Here, i was unable to include an executable because of the way angular works, there is no single executable or html file I can include to just open...
