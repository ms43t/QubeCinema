# Qube Cinema OA Assessment

## Setup Instructions

1. Open the `app.go` file in your project.

2. Initialize the Go module by running the following command in your terminal or command prompt:

    ```bash
    go mod init QubeCinema
    ```

   This command creates a `go.mod` file that tracks the project's dependencies.

3. Run the following command to ensure the `go.mod` file reflects the correct and updated dependencies based on your code:

    ```bash
    go mod tidy
    ```

4. Use the following command to create a `vendor` directory containing all the dependencies required by your project:

    ```bash
    go mod vendor
    ```

   This step ensures that you have a local copy of the dependencies.

5. Build the project by running:

    ```bash
    go build
    ```

   This command compiles the Go code and generates an executable binary file.

6. Finally, run by executing:

    ```bash
    go run app.go
    ```

   This command starts the server and runs your application.

**Note:** Ensure that you are in the correct directory when running these commands. Adjust the commands accordingly if your project structure or file names are different.


## Output

The output of the provided code will be the permissions for each distributor in the specified format. For each distributor, it will print permissions for all cities loaded from the CSV file. The actual output will depend on the data in the CSV file and the permissions specified for each distributor.
