const errorElement = document.querySelector("#table-error");
const searchInput = document.querySelector("#stock-search");
const tableElement = document.querySelector("#stock-table");
const tablePanel = document.querySelector(".table-panel");
const detailView = document.querySelector("#ticker-detail");
const tickerMeta = document.querySelector("#ticker-meta");
const tickerName = document.querySelector("#ticker-name");
const tickerError = document.querySelector("#ticker-error");
const metricSearch = document.querySelector("#metric-search");
const metricSummary = document.querySelector("#metric-summary");

const escapeHtml = (value) =>
    String(value ?? "")
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;");

const formatMetricName = (key) =>
    key
        .replace(/([a-z])([A-Z])/g, "$1 $2")
        .replace(/^./, (letter) => letter.toUpperCase());

const formatValue = (value) =>
    typeof value === "number"
        ? value.toLocaleString("en-US", { maximumFractionDigits: 4 })
        : "—";

const columns = [
    {
        title: "Ticker",
        field: "ticker",
        minWidth: 220,
        widthGrow: 2,
        hozAlign: "left",
        formatter: (cell) => {
            const data = cell.getRow().getData();
            const ticker = escapeHtml(data.ticker);
            return `<a class="ticker-link" href="?ticker=${encodeURIComponent(data.ticker ?? "")}">${ticker}</a> | ${escapeHtml(data.industry)}<br>${escapeHtml(data.name)}`;
        },
    },
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

const svgNS = "http://www.w3.org/2000/svg";

function displayName(name) {
    return formatMetricName(name)
        .replace(/T T M/g, "TTM")
        .replace(/\bE V\b/g, "EV")
        .replace(/\bP B\b/g, "PB")
        .replace(/\bP E\b/g, "PE")
        .replace(/\bP S\b/g, "PS")
        .replace(/\bF C F\b/g, "FCF")
        .replace(/\bR O A\b/g, "ROA")
        .replace(/\bR O E\b/g, "ROE")
        .replace(/\bR O I C\b/g, "ROIC")
        .replace(/\bR O T C\b/g, "ROTC")
        .replace(/\bS G A\b/g, "SGA");
}

function makeSvgElement(tag, attributes) {
    const element = document.createElementNS(svgNS, tag);
    Object.entries(attributes).forEach(([key, value]) => element.setAttribute(key, value));
    return element;
}

function formatPeriod(period, type) {
    const [year, month] = period.split("-").map(Number);
    if (type === "annual") return String(year);
    return `Q${Math.ceil(month / 3)} '${String(year).slice(-2)}`;
}

function makeChart(name, points, type) {
    const styles = getComputedStyle(document.documentElement);
    const colors = {
        positive: styles.getPropertyValue("--teal30").trim(),
        negative: styles.getPropertyValue("--purple40").trim(),
        zero: styles.getPropertyValue("--gray50").trim(),
        baseline: styles.getPropertyValue("--gray70").trim(),
    };
    const width = Math.max(360, points.length * 22);
    const height = 126;
    const left = 8;
    const right = width - 8;
    const top = 18;
    const bottom = 94;
    const values = points.map((point) => point.v);
    const min = Math.min(0, ...values);
    const max = Math.max(0, ...values);
    const span = max - min || 1;
    const zeroY = bottom - ((0 - min) / span) * (bottom - top);
    const slot = (right - left) / points.length;
    const barWidth = Math.max(3, Math.min(16, slot * 0.62));
    const svg = makeSvgElement("svg", {
        class: "metric-chart",
        viewBox: `0 0 ${width} ${height}`,
        role: "img",
        "aria-label": `${displayName(name)} ${type} bar chart`,
    });

    svg.appendChild(makeSvgElement("line", {
        x1: left, y1: zeroY, x2: right, y2: zeroY,
        stroke: colors.baseline, "stroke-width": "1",
    }));

    const labelEvery = Math.max(1, Math.ceil(points.length / 6));
    points.forEach((point, index) => {
        const x = left + slot * index + (slot - barWidth) / 2;
        const valueY = bottom - ((point.v - min) / span) * (bottom - top);
        const rect = makeSvgElement("rect", {
            x, y: Math.min(zeroY, valueY), width: barWidth,
            height: Math.max(1, Math.abs(zeroY - valueY)),
            rx: Math.min(3, barWidth / 3),
            fill: point.v > 0 ? colors.positive : point.v < 0 ? colors.negative : colors.zero,
            opacity: "0.9",
        });
        const tip = makeSvgElement("title", {});
        tip.textContent = `${point.period}: ${formatValue(point.v)}`;
        rect.appendChild(tip);
        svg.appendChild(rect);

        if (index % labelEvery === 0 || index === points.length - 1) {
            const label = makeSvgElement("text", {
                x: left + slot * index + slot / 2,
                y: 116,
                "text-anchor": "middle",
            });
            label.textContent = formatPeriod(point.period, type);
            svg.appendChild(label);
        }
    });

    const latest = makeSvgElement("text", {
        class: "metric-chart-value",
        x: right, y: 11, "text-anchor": "end",
    });
    latest.textContent = `Latest  ${formatValue(points[points.length - 1].v)}`;
    svg.appendChild(latest);
    return svg;
}

function createMetricCard(name, rawPoints, type) {
    const points = rawPoints
        .filter((point) => point && typeof point.period === "string" && Number.isFinite(point.v))
        .slice()
        .sort((a, b) => a.period.localeCompare(b.period))
        .slice(-10);
    if (!points.length) return null;

    const article = document.createElement("article");
    article.className = "metric-card";
    article.dataset.metric = displayName(name).toLocaleLowerCase();

    const head = document.createElement("div");
    head.className = "metric-card-head";
    const title = document.createElement("h4");
    title.textContent = displayName(name);
    const range = document.createElement("span");
    range.className = "metric-range";
    range.textContent = `${formatPeriod(points[0].period, type)} – ${formatPeriod(points[points.length - 1].period, type)}`;
    head.append(title, range);
    article.append(head, makeChart(name, points, type));
    return article;
}

function renderDashboard(series) {
    const counts = {};
    ["annual", "quarterly"].forEach((type) => {
        const grid = document.getElementById(`${type}-grid`);
        grid.replaceChildren();
        Object.entries(series[type] ?? {}).forEach(([name, points]) => {
            const card = createMetricCard(name, points, type);
            if (card) grid.appendChild(card);
        });
        counts[type] = grid.childElementCount;
    });

    function applySearch() {
        const query = metricSearch.value.trim().toLocaleLowerCase();
        let visibleTotal = 0;
        ["annual", "quarterly"].forEach((type) => {
            const grid = document.getElementById(`${type}-grid`);
            let visible = 0;
            Array.from(grid.children).forEach((card) => {
                const show = card.dataset.metric.includes(query);
                card.hidden = !show;
                if (show) visible += 1;
            });
            visibleTotal += visible;
            document.getElementById(`${type}-count`).textContent =
                `${visible} of ${counts[type]} metrics`;
            document.getElementById(`${type}-empty`).hidden = visible !== 0;
        });
        metricSummary.textContent =
            `${visibleTotal} chart${visibleTotal === 1 ? "" : "s"} · ${counts.annual} annual + ${counts.quarterly} quarterly`;
    }

    metricSearch.addEventListener("input", applySearch);
    applySearch();
}

async function loadTickerDetail(ticker) {
    tablePanel.hidden = true;
    detailView.hidden = false;
    document.querySelector(".page-header").hidden = true;
    tickerName.textContent = ticker;
    tickerMeta.textContent = "Loading ticker data…";

    try {
        const [stocksResponse, tickerResponse] = await Promise.all([
            fetch("data/stocks.json"),
            fetch(`data/${encodeURIComponent(ticker)}.json`),
        ]);

        if (!stocksResponse.ok) {
            throw new Error(`Could not load stock information (${stocksResponse.status}).`);
        }
        if (!tickerResponse.ok) {
            throw new Error(`Could not load data for ${ticker} (${tickerResponse.status}).`);
        }

        const stocks = await stocksResponse.json();
        const stock = stocks.find((item) => item.ticker === ticker);
        const data = await tickerResponse.json();

        tickerName.textContent = stock?.name ?? ticker;
        tickerMeta.textContent = stock
            ? `${stock.ticker} | ${stock.industry ?? "Industry unavailable"}`
            : `${ticker} | Ticker data`;
        document.title = `${ticker} | Stock screener | TeeGo`;
        renderDashboard(data);
    } catch (error) {
        tickerError.textContent = error instanceof Error ? error.message : "An unexpected error occurred while loading ticker data.";
        tickerError.hidden = false;
    }
}

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
            resizableColumns: false,
            height: "100%",
            responsiveLayout: "hide",
            placeholder: "No stocks match your search.",
            initialSort: [{ column: "ticker", dir: "asc" }],
        });

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
        errorElement.textContent = error instanceof Error ? error.message : "An unexpected error occurred while loading stocks.";
        errorElement.hidden = false;
    }
}

const requestedTicker = new URLSearchParams(window.location.search).get("ticker");

if (requestedTicker) {
    loadTickerDetail(requestedTicker.toUpperCase());
} else {
    loadStocks();
}