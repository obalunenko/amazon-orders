# amazon-orders

Helper tools to parse amazon order reports files



### Program Structure

The program consists of several parts:

1. **Flag Parsing and File Input**
2. **CSV Parsing**
3. **Order Processing**
4. **Spending Calculation**

### Running the Program

#### Step 1: Prepare CSV Input File

Ensure your CSV input file is correctly formatted and contains the necessary order data. The file should contain headers and data rows corresponding to the order fields.
Reports files could be downloaded from the following URL:
https://www.amazon.com/hz/privacy-central/data-requests/preview.html

#### Step 2: Build and Execute the Program

1. **Compile the Program:**

    ```bash
      make build
    ```

2. **Run the Program:**

    ```bash
       ./bin/order_processor --input=<path/to/your/orders.csv>
    ```

    Replace `<path_to_csv_file>` with the path to your CSV input file.

### Flag Details
- `--input`: The flag to specify the path to the input CSV file. This is a required parameter.

### Example
Here's an example of how to run the program:

```bash
./order_processor --input=orders.csv
```


### Program Execution
Upon running the program with the above command, it will:
1. Read and parse the CSV input file.
2. Extract and categorize orders based on their currencies.
3. Calculate the total spending for each currency.
4. Print the total spend per currency rounded to two decimal places.

### Output
The output will be displayed in the console. Here’s an example of what the output might look like:

```text
Total spend in USD: $1234.56
Total spend in EUR: €789.00
```


### Error Handling
The program includes basic error handling to ensure robustness:
- If the input file is not provided or cannot be opened, a fatal log message will be displayed.
- If there are parsing errors, appropriate error messages will be logged.
