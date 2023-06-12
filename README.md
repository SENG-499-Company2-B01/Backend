# Company 2 | Scheduling Application Backend

Jakobs Scheduler

## Table of Contents

- [Project Description](#project-description)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Examples](#api-examples)
- [Contributing](#contributing)
- [License](#license)

## Project Description

This is a design project to create a course scheduler for the University of Victoria. One of the common problems the university administrators have is the logistic setup of physical resources for courses every academic year. There is a finite number of classrooms, time frames, and specialized equipment that must be assigned. The goal was to build an application based on previous years data that would help automate this process.

One of the main objectives for this project is not only devloping the application but to adhere to the SCRUM process and learn about how to work as part of a larger team. Our backend team has decided to devlop our part using Go as the main language and Python to handle some scripts. We containerize our listening server and Mongo database in docker containers to allow for modularity and future scalability.

## Installation

Before we start, if you do not have Go or Docker on you local computer, you must follow the steps below and install them in order to run this project.

To install Go, please follow these steps:

1. Go to the official Go website by clicking on this link: https://go.dev/doc/install.
2. On the website, locate and click on the blue "Download" button. This will take you to a page where you can select the appropriate download for your Operating System.
3. Choose the download that matches your Operating System and click on it to start the download.
4. Once the download is complete, run the installer and follow the usual steps to install Go.
5. During the installation process, ensure that Go is added to your environmental path. If it doesn't do so automatically, you may need to set it up manually.
6. To verify that Go was installed correctly, open your CMD (command line interface) of choice and enter the command `go --version`. This command should display the installed version of Go if the installation was successful.

To install Docker, please follow these steps:

1. Visit the official Docker website by going to https://www.docker.com/get-started.
2. On the Docker website, you should see a "Get Started with Docker" section. Choose the option that corresponds to your Operating System (e.g., Docker Desktop for Windows, Docker Desktop for Mac).
3. Click on the download link provided for your chosen Operating System.
4. Once the download is complete, run the installer that was downloaded.
5. Follow the on-screen instructions to install Docker on your machine. The installation process may require administrative privileges, so make sure to grant the necessary permissions.
6. During the installation, you may be prompted to configure additional settings or select optional components. Make any desired selections or keep the default options, and proceed with the installation.
7. After the installation is complete, Docker should be available on your computer.
8. To verify that Docker is installed correctly, open a command line interface (e.g., Command Prompt on Windows, Terminal on macOS) and enter the command `docker --version`. This command should display the installed version of Docker if the installation was successful.
9. Optionally, if you are on windows, you may need to install WSL (Windows Subsystem for Linux). To do this:

   - Open PowerShell as an administrator. You can do this by right-clicking on the Start menu, selecting "Windows PowerShell (Admin)."
   - In the PowerShell window, run the following command to enable the Windows Subsystem for Linux feature:

     `wsl --install`

     This command will automatically enable the necessary components and install a WSL-compatible Linux distribution from the Microsoft Store.

   - Wait for the installation process to complete. It may take a few minutes as it downloads and sets up the Linux distribution.
   - Once the installation is finished, you will be prompted to create a new user account and set a password for the Linux distribution. Follow the on-screen instructions to complete this step.
   - After creating the user account, WSL is ready to use.

By following these steps, you will successfully install Go and Docker on your machine.

Now we can move onto installing the project on your machine in order to configure and run it.

To start, clone the repository onto your machine. If you are unsure on how to do this:

1. Continue to the root page of this repository.
2. Click on the blue '<> Code' button near the top and copy the HTTPS link.
3. Open your chosen command line interface and type the command

   `git clone https://github.com/SENG-499-Company2-B01/Backend.git`

After is completes, you should be able to open the project in your prefered development environment.

## Configuration

Configuration for this project is very simple!

Open the project in your developement environment and create a new file in the root folder called `.env`. In there you must follow this exact format, but feel free to change the values for MONGO_USERNAME, MONGO_PASSWORD, and JWT_SECRET.

```
MONGO_ADDRESS=10.9.0.3
MONGO_PORT=27017
MONGO_USERNAME=admin
MONGO_PASSWORD=admin
JWT_SECRET=secret
```

## Usage

The usage of our project is also very simply!

Firstly, you need to start docker for the containers the project will run in, use the following steps:

1. Open Docker Desktop on your computer. This application manages the Docker environment and allows you to control containers.

   - If you are using Windows or macOS, you can find the Docker Desktop icon in the system tray or the applications menu. Click on it to open Docker Desktop.
   - If you are using Linux, open a terminal and run the command to start the docker service on your machine.

     ```
     sudo service docker start
     ```

2. To verify that docker is running, you can use the following command in your CLI. If it returns an error, docker is either not running or was not set up correctly.

   ```
   docker info
   ```

3. After confirming that Docker is running, navigate to the directory where your project's docker-compose.yml file is located. This file defines the containers and their configurations.
4. To start the containers defined in your docker-compose.yml file, run the following command:

   ```
   docker-compose up
   ```

   Optionally, you can add the -d flag at the end of the command to detach the containers from the CLI. This allows them to run in the background while you continue to use the CLI for other tasks.

   ```
   docker-compose up -d
   ```

   This command will spin up the containers and their dependencies, as specified in the docker-compose.yml file.

You are now ready to interact with the server. This backend part of the project is strictly an interface to the outside to allow for the other aspects of the project to communicate.

To interact with one of the server endpoints, you must send a request like any other API. Our services include:

- Signin
- Users
- Classrooms
- Courses
- Schedules

## API Examples

Some examples of the requests are as follows:

### Endpoint: Retrieve Specific User

**Endpoint:** `http://localhost:8000/users/{id}`

**Method:** `GET`

**Description:** `Retrieves a specific user from the database if they exist.`

**Example Request:**

```
http
GET /users/test?test=example HTTP/1.1
Host: localhost:8000
Accept: application/json
```

**Example Request:**

```
HTTP/1.1 200 OK
Content-Type: application/json

{
    "username":"test",
    "email":"test@telus.net",
    "password":"password",
    "firstname":"John",
    "lastname":"Doe",
    "preferences":{
        "key1":"value1",
        "key2":"value2"
    },
    "qualifications":["qualification1","qualification2"]
}
```

### Endpoint: Create a new course

**Endpoint:** `http://localhost:8000/courses`

**Method:** `POST`

**Description:** Creates a new course.

**Request Body:**

- `shorthand` (string, required): The shorthand code for the course.
- `name` (string, required): The name of the course.
- `equipments` (array of strings, optional): The list of equipment required for the course.
- `prerequisites` (array of strings, optional): The list of prerequisite courses for the course.

**Example Request:**

```
http
POST /courses HTTP/1.1
Host: localhost:8000
Content-Type: application/json

{
  "shorthand": "SENG499",
  "name": "Design Project II",
  "equipments": ["projector", "camera"],
  "prerequisites": ["SENG360", "SENG210"]
}
```

**Example Request:**

```
HTTP/1.1 200 OK
Content-Type: application/json

{
    "shorthand":"SENG499",
    "name":"Design Project II",
    "equipment":null,
    "prerequisites":["SENG360","SENG 210"]
}
```

For more endpoint examples, please contact someone from our backend team, or refer to the Company SRS document.

## Contributing

As this is currently a private project without external help, feel free to reach out to our backend team if you want any changes or help and we will get back to you as soon as possible.

For developers:

- For bug reports clearly outline what the bug is, how to reproduce it and provide screenshots or a video.
- For feature requests, clearly outline what the requirements for the feature are and why you think this feature is needed.
- For pull requests, clearly outline what feature this is for (ticket number), what changes were made, and how to test is with optionally screenshots or a video.

## License

This project is licensed under the [MIT License](LICENSE).

[Click here to view the license file](LICENSE) and review the terms and conditions of the MIT License.

## Contact

Contact the backend team for Company 2 :)

**README.md last Updated 2023-06-12**
