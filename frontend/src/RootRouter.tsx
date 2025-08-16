import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import App from './App';
import SearchPage from './SearchPage';

export default function RootRouter() {
  return (
    <Router>
      <Routes>
  <Route path="/" element={<App />} />
  <Route path="/search" element={<SearchPage />} />
      </Routes>
    </Router>
  );
}
