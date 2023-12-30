## 👻 Tenebra

### Introduction
Tenebra is a GoLang project inspired by Projectdiscovery's Chaos, designed to bring a local version of Chaos DB to your environment. It allows you to fetch datasets from the Chaos portal, download, clean, and merge them into a single text file. Additionally, Tenebra sets up a local API endpoint, enabling you to search for DNS/URLs within the cached dataset.

#### Features

- **Data Fetching:** Fetches content from the Chaos portal via the JSON endpoint.
- **Dataset Download:** Downloads and unzips all datasets pointed to in the JSON items.
- **Cleanup:** Removes invalid DNS/URLs from the datasets.
- **Merge:** Merges cleaned datasets into a single "tenebra.txt" file.
- **Local API:** Runs a local API endpoint for searching DNS/URLs in the cached "tenebra.txt."

## Getting Started

**Prerequisites:**
- Go (Golang) installed on your system.
- Internet connection to fetch initial data from the Chaos portal.

**Installation:**
1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/tenebra.git
   cd tenebra
2. Install dependencies
   ```bash
   go mod tidy
3. Build
   ```bash
   go build
4. Run
   ```bash
   ./tenebra

## API Usage

Once Tenebra is running, you can use the local API to search for DNS/URLs.
```
Endpoint: localhost:1991/search?query=
Method: GET
Query Parameter: query (e.g., /search?query=example.com)
```

#### Contributing
Contributions are welcome! Feel free to open issues or submit pull requests.

#### License
This project is licensed under the MIT License - see the LICENSE file for details.

#### Acknowledgments
This project was inspired by Projectdiscovery's Chaos.