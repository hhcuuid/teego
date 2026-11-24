const countElement = document.querySelector("#stock-count");
const errorElement = document.querySelector("#table-error");
const searchInput = document.querySelector("#stock-search");
const tableElement = document.querySelector("#stock-table");

const columns = [
    { title: "Ticker", field: "ticker", minWidth: 90, width: 110, hozAlign: "left" },
    { title: "Company", field: "name", minWidth: 200, widthGrow: 2 },
    { title: "Industry", field: "industry", minWidth: 140, widthGrow: 1 },
    {
        title: "Market cap (B)",
        field: "mcapbillion",
        sorter: "number",
        minWidth: 135,
        hozAlign: "right",
        formatter: "money",
        formatterParams: { decimalPlaces: 2, symbol: "$", symbolAfter: false },
    },
    { title: "Beta", field: "beta", sorter: "number", minWidth: 90, hozAlign: "right", formatter: "money", formatterParams: { decimalPlaces: 2, symbol: "" } },
    {
        title: "Dividend",
        field: "dividend",
        sorter: "number",
        minWidth: 110,
        hozAlign: "right",
        formatter: "money",
        formatterParams: { decimalPlaces: 2, symbol: "$", symbolAfter: false },
    },
    { title: "P/E (TTM)", field: "pettm", sorter: "number", minWidth: 110, hozAlign: "right", formatter: "money", formatterParams: { decimalPlaces: 2, symbol: "" } },
    {
        title: "10d vs 90d vol. (%)",
        field: "vol10d90d",
        sorter: "number",
        minWidth: 170,
        hozAlign: "right",
        formatter: "money",
        formatterParams: { decimalPlaces: 2, symbol: "%", symbolAfter: true },
    },
];

async function loadStocks() {
    try {
        const response = await fetch("data/stocks.json");
        if (!response.ok) {
            throw new Error(`Could not load stocks (${response.status}).`);
        }

        const stocks = await response.json();
        if (!Array.isArray(stocks)) {
            throw new Error("Stock data must be a JSON array.");
        }

        const table = new Tabulator(tableElement, {
            data: stocks,
            columns,
            layout: "fitColumns",
            height: "100%",
            responsiveLayout: "hide",
            placeholder: "No stocks match your search.",
            initialSort: [{ column: "ticker", dir: "asc" }],
        });

        const updateCount = (rows) => {
            const visible = rows.length;
            const total = stocks.length;
            countElement.textContent = visible === total
                ? `${total} ${total === 1 ? "stock" : "stocks"}`
                : `${visible} of ${total} stocks`;
        };

        table.on("dataFiltered", (_filters, rows) => updateCount(rows));
        updateCount(stocks);

        searchInput.addEventListener("input", () => {
            const query = searchInput.value.trim().toLowerCase();
            if (!query) {
                table.clearFilter();
                return;
            }

            table.setFilter((data) =>
                [data.ticker, data.name, data.industry]
                    .some((value) => String(value ?? "").toLowerCase().includes(query)),
            );
        });
    } catch (error) {
        countElement.textContent = "Stocks unavailable";
        errorElement.textContent = error instanceof Error ? error.message : "An unexpected error occurred while loading stocks.";
        errorElement.hidden = false;
    }
}

loadStocks();