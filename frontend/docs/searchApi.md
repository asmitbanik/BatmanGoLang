# searchApi.ts

## Purpose
Handles API calls for code search from the frontend.

## Features
- Calls `/search` endpoint with query, filters, pagination.
- Maps backend results to frontend format.
- (Optionally) can be extended for GitHub search.

## Usage
- Import and use `searchCode(query, source, filters, limit, offset)` in your React components.

---

*Updated: Now documents backend integration and result mapping.*
