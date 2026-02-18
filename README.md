# gate-keeper

<!-- [![Go Coverage](https://github.com/matttm/gate-keeper/wiki/coverage.svg)](https://raw.githack.com/wiki/matttm/gate-keeper/coverage.html) -->

## What is gate-keeper?

Gate-keeper is a tool that helps you manage application gates in a database. It allows you to adjust gate open and close times so you can test your application at different stages of an application cycle.

### What does it do?

Think of your application cycle as a timeline with different phases (gates). Gate-keeper lets you move to any point in that timeline by automatically updating the dates in your database.

**Example Timeline:**
```
       aaaaa       bbbbb       ccccc       ddddd       eeeee     
------------------------------------------------------------------------
```
This represents an application cycle with five gates (phases a through e).

**How it works:**
1. You select a year (e.g., 2025)
2. You select a gate (e.g., gate c)
3. You choose where you want to be positioned:
   - **Before** - Places you just before the gate opens
   - **Inside** - Places you during the gate's open period
   - **After** - Places you just after the gate closes

The program then updates all gate dates in the database so the current time appears to be at your chosen position.

**Visual Examples:**

Before gate c:
```
       aaaaa       bbbbb  me   ccccc       ddddd       eeeee     
------------------------------------------------------------------------
```

Inside gate c:
```
                                me
       aaaaa       bbbbb       ccccc       ddddd       eeeee     
------------------------------------------------------------------------
```

After gate c:
```
       aaaaa       bbbbb       ccccc   me  ddddd       eeeee     
------------------------------------------------------------------------
```

## Quick Start (Recommended for Testers)

### What You'll Need

1. **A configuration file** - This tells the program how to connect to your database
2. **The gate-keeper program** - Either download a pre-built version or build it yourself

### Step 1: Download the Program

**Option A: Download Pre-built Version (Easiest)**

1. Go to the [Releases section](../../releases) of this repository
2. Download the file that matches your system:
   - Windows: Look for a file ending in `.exe`
   - Mac: Look for files labeled `darwin` (choose `arm64` for newer Macs with M1/M2 chips, or `amd64` for older Intel Macs)
   - Linux: Choose based on your system architecture (most likely `amd64`)
3. Save the downloaded file to a folder on your computer

**Option B: Build from Source (For Developers)**

If you have Go installed and want to build from source:

1. Open a terminal/command prompt
2. Navigate to where you want to place the program
3. Run these commands:
   ```bash
   go mod download
   go build
   ```
4. This creates a `gate-keeper` program file in your current folder

### Step 2: Create Your Configuration File

You need to create a file called `config.json` in the **same folder** as the gate-keeper program.

**How to create config.json:**

1. Open a text editor (Notepad on Windows, TextEdit on Mac, or any code editor)
2. Copy the template below
3. Replace the placeholder values with your actual database information
4. Save the file as `config.json` in the same folder as gate-keeper

**Template:**
```json
{
	"Credentials": {
		"User": "your_database_username",
		"Pass": "your_database_password",
		"Host": "your_database_server_address",
		"Port": "3306"
	},
	"GateConfig": {
		"Dbname": "your_database_name",
		"TableName": "your_gates_table_name",
		"GateNameKey": "gate_name_column",
		"GateYearKey": "year_column",
		"GateOrderKey": "order_column",
		"GateIsApplicableFlag": "active_flag_column",
		"StartKey": "start_date_column",
		"EndKey": "end_date_column"
	},
	"EnablePprof": false
}
```
	},
	"GateConfig": {
		"Dbname": "your_database_name",
		"TableName": "your_gates_table_name",
		"GateNameKey": "gate_name_column",
		"GateYearKey": "year_column",
		"GateOrderKey": "order_column",
		"GateIsApplicableFlag": "active_flag_column",
		"StartKey": "start_date_column",
		"EndKey": "end_date_column"
	}
}
```

**What each field means:**

**Credentials Section** (Database connection details):
- `User`: Your database username
- `Pass`: Your database password
- `Host`: The address of your database server (e.g., "localhost" or "db.yourcompany.com")
- `Port`: The database port (usually "3306" for MySQL)

**GateConfig Section** (Tells the program where to find gate information in your database):
- `Dbname`: The name of your database
- `TableName`: The name of the table that contains gate information
- `GateNameKey`: The column name that stores the gate name (e.g., "code" or "gate_name")
- `GateYearKey`: The column name that stores the year (e.g., "year")
- `GateOrderKey`: The column name that defines the order of gates (e.g., "order" or "sequence")
- `GateIsApplicableFlag`: The column name that indicates if a gate is active (e.g., "active" or "enabled")
- `StartKey`: The column name for when a gate opens (e.g., "start" or "start_date")
- `EndKey`: The column name for when a gate closes (e.g., "end" or "end_date")

**Performance Profiling Section**:
- `EnablePprof`: Set to `true` to enable the pprof profiling server on port 8080 (used for performance analysis and flamegraphs). Should be `false` for normal use.

**Example with real values:**
```json
{
	"Credentials": {
		"User": "test_user",
		"Pass": "myPassword123",
		"Host": "localhost",
		"Port": "3306"
	},
	"GateConfig": {
		"Dbname": "application_db",
		"TableName": "application_gates",
		"GateNameKey": "gate_code",
		"GateYearKey": "cycle_year",
		"GateOrderKey": "sequence_number",
		"GateIsApplicableFlag": "is_active",
		"StartKey": "open_date",
		"EndKey": "close_date"
	},
	"EnablePprof": false
}
```

### Step 3: Run the Program

**On Windows:**
1. Double-click the `gate-keeper.exe` file
2. If Windows shows a security warning, click "More info" then "Run anyway"

**On Mac/Linux:**
1. Open a terminal
2. Navigate to the folder containing gate-keeper
3. Run: `./gate-keeper`
4. If you get a permission error, first run: `chmod +x gate-keeper`

### Step 4: Using the Interface

When the program opens, you'll see a window with three dropdown menus and a button:

<img width="519" alt="Screenshot 2025-06-08 at 2 08 19 PM" src="https://github.com/user-attachments/assets/9b1a2f84-7216-49d6-8683-57eb815d72c7" />

**To set gates:**

1. **Select a year** - Choose the application cycle year from the first dropdown
2. **Select a gate** - Choose which gate you want to position yourself relative to
3. **Position relative to gate** - Choose where you want to be:
   - "Before" - Before the gate opens
   - "Inside" - While the gate is open
   - "After" - After the gate closes
4. Click the **"Set Gates"** button to update the database

**Understanding the Gate Status Table:**

Once you select a year, you'll see a table showing all gates and their current status:

<img width="509" alt="Screenshot 2025-06-08 at 2 14 05 PM" src="https://github.com/user-attachments/assets/cdcb735b-8ea8-4abf-bdcf-401d3399dad8" />
<img width="504" alt="Screenshot 2025-06-08 at 2 10 06 PM" src="https://github.com/user-attachments/assets/322c8f5c-90fa-4a20-868a-34543570fb74" />

**Status Colors:**
- **Green** - Gate is currently open
- **Red** - Gate is closed (either not started yet or already ended)
- **Yellow** - Warning: Gates are out of order in the database

**Status Labels:**
- **Past** - The gate's time period has already ended
- **Open** - The gate is currently open
- **Future** - The gate hasn't opened yet

The table updates automatically every second, so you can see status changes in real-time.

**Gate Health Check:**

If all gates turn yellow, this means the gates are out of order in your database:

<img width="522" alt="Screenshot 2025-06-08 at 2 08 52 PM" src="https://github.com/user-attachments/assets/49c5c40d-c8be-4901-98a3-ccc27aa4f8d7" />

This typically means the dates don't follow a logical sequence (e.g., gate B starts before gate A ends).

## Troubleshooting

### The program won't start

**Problem:** Double-clicking the program does nothing, or you see an error about missing files.

**Solution:**
- Make sure `config.json` is in the same folder as the gate-keeper program
- Check that your config.json file is valid JSON (you can use an online JSON validator)
- On Mac/Linux, make sure the file has execute permissions: `chmod +x gate-keeper`

### "Error getting executable path" or "Error opening config file"

**Problem:** The program can't find the config.json file.

**Solution:**
- Verify that `config.json` is in the exact same folder as the gate-keeper program
- Check the file name is exactly `config.json` (not `config.json.txt` or any other variation)
- Make sure the file has the correct JSON format

### Can't connect to database

**Problem:** The program starts but shows errors about database connection.

**Solution:**
- Verify your database credentials in config.json are correct
- Make sure the database server is running and accessible
- Check that the Host and Port values are correct
- Ensure your database user has permission to read and write to the gates table

### "All selections are required" popup appears

**Problem:** When you click "Set Gates", a popup says all selections are required.

**Solution:**
- Make sure you've selected a value from all three dropdowns (year, gate, and position)
- Try selecting each dropdown again to ensure the values are registered

### Gates appear out of order (all yellow)

**Problem:** All gates in the table are yellow instead of green/red.

**Solution:**
- This means the dates in your database are not in sequential order
- This might be the expected state if you're testing edge cases
- To fix: Use gate-keeper to set gates to a normal position, which should correct the ordering

### Wrong columns or data

**Problem:** The program shows incorrect gate names or data.

**Solution:**
- Double-check the column names in your `GateConfig` section of config.json
- Make sure each field matches your actual database table column names exactly (including case sensitivity)
- Verify you're pointing to the correct database and table

## For Developers

### Building from Source

Requirements:
- Go 1.16 or later
- MySQL database

Build instructions:
```bash
go mod download
go build
```

### Creating Binaries for Multiple Platforms

The included `build.sh` script can build binaries for different Linux architectures:

```bash
./build.sh
```

This creates binaries for: amd64, arm64, ppc64le, ppc64, and s390x.

For other platforms, use:
```bash
# For Windows
env GOOS=windows GOARCH=amd64 go build -o gate-keeper.exe

# For Mac (Intel)
env GOOS=darwin GOARCH=amd64 go build -o gate-keeper-mac-intel

# For Mac (M1/M2)
env GOOS=darwin GOARCH=arm64 go build -o gate-keeper-mac-arm
```

## Technical Details

### How it works

The program updates all gates for a selected year by:
1. Querying the database for all gates in the specified year
2. Calculating the time offsets needed to place the current time at the desired position
3. Updating the start and end datetime columns for all gates in the cycle

The position is calculated relative to the selected gate:
- **Before**: Current time is set to be slightly before the gate's start time
- **Inside**: Current time is set to be within the gate's start and end time
- **After**: Current time is set to be slightly after the gate's end time

All other gates in the cycle are adjusted proportionally to maintain their relative durations and spacing.

## Changelog
- 06/08/2025 - gate status table with a gate health check

## Authors

-   Matt Maloney : matttm

## Contribute

If you want to contribute, just send me a message.
