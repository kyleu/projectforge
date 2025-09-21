# Search

A comprehensive search system that provides full-text search capabilities across your application's data. The search component includes a user-friendly search bar interface and a powerful backend search engine that can index and search through various data types.

- **Programmable Search**: Search across multiple fields and data types using your own code
- **Configurable**: Enable/disable search on specific pages
- **Performance Optimized**: Efficient indexing and querying
- **Type-Safe**: Go-based search definitions with proper typing
- **Export Integration**: Automatic search generation for exported models

## Prerequisites

The search functionality requires the `search` module to be enabled in your Project Forge application. If you're using model exports, the `export` module should also be enabled for automatic search function generation.

## Basic Usage

### Enabling Search

Search is automatically included in your pages when the search module is enabled. The search bar appears in the standard page layout.

### Disabling Search on Specific Pages

To disable the search bar on specific pages, set the `SearchPath` to `"-"` in your controller action:

```go
func (s *Service) HandleSpecialPage(w http.ResponseWriter, r *http.Request, ps *cutil.PageState) {
    // Disable search bar on this page
    ps.SearchPath = "-"

    // Your page logic here
    return s.Render(w, r, ps, "special_page", "Special Page")
}
```

### Custom Search Path

You can customize where search results are handled by setting a custom search path:

```go
func (s *Service) HandleProductPage(w http.ResponseWriter, r *http.Request, ps *cutil.PageState) {
    // Use custom search path for product-specific search
    ps.SearchPath = "/products/search"

    return s.Render(w, r, ps, "products", "Products")
}
```

## Search Engine Architecture

The search system is built around the components in `/app/lib/search`:

### Core Components

1. **Search Interface**: Defines the contract for search operations
2. **Result Types**: Structured result objects for different content types
3. **Indexing**: Content indexing for efficient searching
4. **Query Processing**: Search query parsing and execution

### Search Types

The search system defines several types for handling search operations:

```go
// Example search result structure
type SearchResult struct {
    ID          string                 `json:"id"`
    Type        string                 `json:"type"`
    Title       string                 `json:"title"`
    Summary     string                 `json:"summary"`
    URL         string                 `json:"url"`
    Score       float64               `json:"score"`
    Metadata    map[string]interface{} `json:"metadata"`
    Highlights  []string              `json:"highlights"`
}

// Search query parameters
type SearchParams struct {
    Query    string   `json:"query"`
    Types    []string `json:"types,omitempty"`
    Limit    int      `json:"limit,omitzero"`
    Offset   int      `json:"offset,omitzero"`
    Filters  map[string]interface{} `json:"filters,omitzero"`
}
```

## Model-Based Search

### Automatic Search Generation

When the `export` module is enabled and models have `search = true` configured, search functions are automatically generated in `app/lib/search/generated.go`.

#### Model Configuration

```go
// In your model definition
type User struct {
    ID       string `json:"id" search:"true"`
    Name     string `json:"name" search:"true"`
    Email    string `json:"email" search:"true"`
    Bio      string `json:"bio" search:"true"`
    Active   bool   `json:"active"`
    Created  time.Time `json:"created"`
}
```

#### Generated Search Functions

The system automatically generates search functions like:

```go
// Generated search function for User model
func SearchUsers(ctx context.Context, db *sql.DB, query string, params *SearchParams) ([]*SearchResult, error) {
    // Implementation generated based on model definition
    // Searches across Name, Email, and Bio fields
    // Returns structured SearchResult objects
}
```

### Custom Search Implementation

You can also implement custom search functions for specific needs:

```go
// Custom search function
func SearchProducts(ctx context.Context, db *sql.DB, query string, params *SearchParams) ([]*SearchResult, error) {
    var results []*SearchResult

    // Build search query
    sqlQuery := `
        SELECT id, name, description, price, category
        FROM products
        WHERE (name ILIKE $1 OR description ILIKE $1)
        AND active = true
        ORDER BY
            CASE WHEN name ILIKE $1 THEN 1 ELSE 2 END,
            name
        LIMIT $2 OFFSET $3
    `

    searchTerm := "%" + query + "%"
    rows, err := db.QueryContext(ctx, sqlQuery, searchTerm, params.Limit, params.Offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var id, name, description, category string
        var price float64

        err := rows.Scan(&id, &name, &description, &price, &category)
        if err != nil {
            continue
        }

        result := &SearchResult{
            ID:      id,
            Type:    "product",
            Title:   name,
            Summary: description,
            URL:     "/products/" + id,
            Metadata: map[string]interface{}{
                "price":    price,
                "category": category,
            },
        }

        results = append(results, result)
    }

    return results, nil
}
```

## Search Interface Implementation

### Search Handler

```go
func (s *Service) HandleSearch(w http.ResponseWriter, r *http.Request, ps *cutil.PageState) {
    query := r.URL.Query().Get("q")
    if query == "" {
        // Handle empty search
        return s.Render(w, r, ps, "search", "Search")
    }

    // Perform search across different content types
    results, err := s.performSearch(r.Context(), query)
    if err != nil {
        // Handle search error
        ps.Logger.Error("search failed", "error", err, "query", query)
        return controller.FlashAndRedir(false, "Search failed. Please try again.", "/search", w, ps)
    }

    // Render search results
    ps.Data = map[string]interface{}{
        "query":   query,
        "results": results,
        "count":   len(results),
    }

    return s.Render(w, r, ps, "search_results", "Search Results")
}

func (s *Service) performSearch(ctx context.Context, query string) ([]*SearchResult, error) {
    var allResults []*SearchResult

    params := &SearchParams{
        Query: query,
        Limit: 50,
    }

    // Search users
    userResults, err := SearchUsers(ctx, s.db, query, params)
    if err != nil {
        return nil, err
    }
    allResults = append(allResults, userResults...)

    // Search products
    productResults, err := SearchProducts(ctx, s.db, query, params)
    if err != nil {
        return nil, err
    }
    allResults = append(allResults, productResults...)

    // Sort results by relevance score
    sort.Slice(allResults, func(i, j int) bool {
        return allResults[i].Score > allResults[j].Score
    })

    return allResults, nil
}
```

### Search Results Template

```html
<div class="search-results">
  <div class="search-header">
    <h2>Search Results</h2>
    <p>Found {%d len(results) %} results for "{%s query %}"</p>
  </div>

  {%- if len(results) > 0 -%}
  <div class="results-list">
    {%- for _, result := range results -%}
    <div class="result-item" data-type="{%s result.Type %}">
      <div class="result-header">
        <h3 class="result-title">
          <a href="{%s result.URL %}">{%s result.Title %}</a>
        </h3>
        <span class="result-type">{%s result.Type %}</span>
      </div>

      {%- if result.Summary != "" -%}
      <p class="result-summary">{%s result.Summary %}</p>
      {%- endif -%}

      {%- if len(result.Highlights) > 0 -%}
      <div class="result-highlights">
        {%- for _, highlight := range result.Highlights -%}
        <span class="highlight">{%s highlight %}</span>
        {%- endfor -%}
      </div>
      {%- endif -%}

      <div class="result-meta">
        <a href="{%s result.URL %}" class="result-link">View Details</a>
        {%- if result.Score > 0 -%}
        <span class="result-score">Score: {%0.2f result.Score %}</span>
        {%- endif -%}
      </div>
    </div>
    {%- endfor -%}
  </div>
  {%- else -%}
  <div class="no-results">
    <h3>No results found</h3>
    <p>Try adjusting your search terms or browse our content directly.</p>
    <div class="search-suggestions">
      <h4>Suggestions:</h4>
      <ul>
        <li>Check your spelling</li>
        <li>Try more general terms</li>
        <li>Use fewer keywords</li>
      </ul>
    </div>
  </div>
  {%- endif -%}
</div>
```

## Advanced Search Features

### Filtered Search

```go
func (s *Service) HandleAdvancedSearch(w http.ResponseWriter, r *http.Request, ps *cutil.PageState) {
    query := r.URL.Query().Get("q")
    contentType := r.URL.Query().Get("type")
    category := r.URL.Query().Get("category")

    params := &SearchParams{
        Query: query,
        Types: []string{contentType},
        Filters: map[string]interface{}{
            "category": category,
        },
        Limit: 20,
    }

    results, err := s.performFilteredSearch(r.Context(), params)
    if err != nil {
        return controller.FlashAndRedir(false, "Search failed", "/search", w, ps)
    }

    ps.Data = map[string]interface{}{
        "query":   query,
        "results": results,
        "filters": params.Filters,
    }

    return s.Render(w, r, ps, "advanced_search", "Advanced Search")
}
```

### Search with Autocomplete

```html
<form class="search-form" action="/search" method="get">
  <div class="search-input-group">
    <input
      type="text"
      name="q"
      class="search-input"
      placeholder="Search..."
      value="{%s query %}"
      autocomplete="off"
      data-autocomplete="/api/search/suggestions"
    >
    <button type="submit" class="search-button">
      <span class="search-icon">🔍</span>
    </button>
  </div>

  <div class="search-suggestions" id="search-suggestions" style="display: none;">
    <!-- Autocomplete suggestions populated via JavaScript -->
  </div>
</form>
```

### Search Analytics

```go
// Track search queries for analytics
func (s *Service) trackSearch(ctx context.Context, query string, resultCount int, userID string) {
    searchLog := &SearchLog{
        Query:       query,
        ResultCount: resultCount,
        UserID:      userID,
        Timestamp:   time.Now(),
        IP:          getClientIP(ctx),
    }

    // Store search analytics
    err := s.db.CreateSearchLog(ctx, searchLog)
    if err != nil {
        s.logger.Error("failed to log search", "error", err)
    }
}
```

## Best Practices

### Performance Optimization
1. **Indexing**: Create appropriate database indexes for searchable fields
2. **Caching**: Cache frequent search results
3. **Pagination**: Implement pagination for large result sets
4. **Query Optimization**: Use efficient SQL queries with proper LIMIT/OFFSET

### User Experience
1. **Fast Response**: Keep search response times under 200ms when possible
2. **Relevant Results**: Order results by relevance and recency
3. **Clear Feedback**: Show search progress and result counts
4. **Error Handling**: Provide helpful messages for failed searches

### Search Quality
1. **Fuzzy Matching**: Handle typos and similar terms
2. **Stemming**: Match word variations (run, running, ran)
3. **Synonyms**: Include related terms in search results
4. **Stop Words**: Filter out common words that don't add value

### Security
1. **Input Sanitization**: Clean search queries to prevent injection attacks
2. **Rate Limiting**: Prevent search abuse with rate limiting
3. **Access Control**: Respect user permissions in search results
4. **Logging**: Log search activities for security monitoring

## Common Use Cases

### Content Management
- Search through articles, pages, and media
- Find content by title, tags, or full text
- Filter by publication status or date

### E-commerce
- Product search by name, description, or SKU
- Filter by category, price range, or availability
- Search customer orders and information

### User Management
- Find users by name, email, or role
- Search user-generated content
- Administrative user lookup

### Documentation
- Search help articles and documentation
- Find API endpoints or code examples
- Search through comments and discussions

## Troubleshooting

### Search Not Appearing
- Verify the `search` module is enabled
- Check that `ps.SearchPath` is not set to `"-"`
- Ensure search templates are properly included

### Poor Search Results
- Review database indexes on searchable fields
- Check search query logic and scoring
- Verify data is properly indexed

### Performance Issues
- Add database indexes for search queries
- Implement result caching
- Consider using dedicated search engines (Elasticsearch, etc.)

### Integration Problems
- Ensure `export` module is enabled for model search
- Check that models have `search = true` configured
- Verify generated search functions are properly wired
