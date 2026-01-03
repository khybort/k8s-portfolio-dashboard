import React from 'react';
import ReactDOM from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import AdminApp from './admin/App';
import PublicApp from './public/PublicApp';
import './index.css';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

// Determine which app to render based on domain and path
const AppRouter = () => {
  const hostname = window.location.hostname;
  const path = window.location.pathname;
  
  // If accessing via admin.portfolio.local, always show admin app
  if (hostname === 'admin.portfolio.local' || hostname.startsWith('admin.')) {
    return <AdminApp />;
  }
  
  // If path starts with /admin, show admin app
  if (path.startsWith('/admin')) {
    return <AdminApp />;
  }
  
  // Otherwise, show public app
  return <PublicApp />;
};

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <AppRouter />
    </QueryClientProvider>
  </React.StrictMode>
);
