package newsapi

const (
	NewsAPIBASEURL         = "https://newsapi.org"
	NewsAPIURLEverything   = NewsAPIBASEURL + "/v2/everything"
	NewsAPIURLTopHeadlines = NewsAPIBASEURL + "/v2/top-headlines"
	NewsAPIURLSources      = NewsAPIURLTopHeadlines + "/sources"
)

type NewsAPIEverythingReq struct {
	Q              string `json:"q,omitempty"`
	Searchin       string `json:"searchIn,omitempty"`
	PageSize       int    `json:"pageSize,omitempty"`
	Page           int    `json:"page,omitempty"`
	SortBy         string `json:"sortBy,omitempty"`         // options: relevancy, popularity, publishedAt
	From           string `json:"from,omitempty"`           // YYYY-MM-DD
	To             string `json:"to,omitempty"`             // YYYY-MM-DD
	Domains        string `json:"domains,omitempty"`        // comma separated list of domains to include in the results
	ExcludeDomains string `json:"excludeDomains,omitempty"` // comma separated list of domains to exclude from the results
	Language       string `json:"language,omitempty"`       // ISO 639-1 code of the language to restrict the results to, Possible options: ar, de, en, es, fr, he, it, nl, no, pt, ru, sv, ud, zh.
	// Sources        string `json:"sources,omitempty"`        // comma separated list of identifiers for the news sources or blogs to include in the results
}

type NewsAPITopHeadlinesReq struct {
	Country  string `json:"country,omitempty"`  // ISO 3166-1 code of the country to restrict the results to BUT for this in Docs just have can use one country code -> "us"
	Category string `json:"category,omitempty"` // options: business, entertainment, general, health, science, sports, technology
	// Sources  string `json:"sources,omitempty"`  // comma separated list of identifiers for the news sources or blogs to include in the results
	Q        string `json:"q,omitempty"`        // keywords or a phrase to search for in the article title and body
	PageSize int    `json:"pageSize,omitempty"` // number of results to return per page (1-100)
	Page     int    `json:"page,omitempty"`     // page number to retrieve (1-indexed)
}

type ArticleSource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	Source      ArticleSource `json:"source"`
	Author      string        `json:"author"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Url         string        `json:"url"`
	UrlToImage  string        `json:"urlToImage"`
	PublishedAt string        `json:"publishedAt"` // YYYY-MM-DDTHH:MM:SSZ
	Content     string        `json:"content"`
}

type NewsAPIResponse struct {
	Status string `json:"status"` // "ok" or "error"
	// for success response
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
	// for error response
	Code    string `json:"code"`
	Message string `json:"message"`
}
