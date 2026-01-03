import { BrowserRouter, Routes, Route } from 'react-router-dom';
import HomePage from './pages/HomePage';
import ArticleDetailPage from './pages/ArticleDetailPage';

function PublicApp() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/articles/:slug" element={<ArticleDetailPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default PublicApp;

