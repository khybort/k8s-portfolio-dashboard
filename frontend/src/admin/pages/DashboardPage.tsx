import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { motion } from 'framer-motion';
import { 
  FileText, 
  FolderKanban, 
  User, 
  TrendingUp,
  Eye,
  Calendar,
  ArrowRight
} from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { articleService } from '../../services/api/articleService';
import { projectService } from '../../services/api/projectService';
import { portfolioService } from '../../services/api/portfolioService';

const DashboardPage: React.FC = () => {
  const navigate = useNavigate();

  const { data: articles } = useQuery({
    queryKey: ['articles', 1],
    queryFn: () => articleService.getArticles(1, 5),
  });

  const { data: projects } = useQuery({
    queryKey: ['projects', 1],
    queryFn: () => projectService.getProjects(1, 5),
  });

  const { data: portfolio } = useQuery({
    queryKey: ['portfolio'],
    queryFn: () => portfolioService.getPortfolio(),
  });

  const stats = [
    {
      name: 'Total Articles',
      value: articles?.total || 0,
      icon: FileText,
      color: 'bg-blue-500',
      href: '/admin/articles',
    },
    {
      name: 'Total Projects',
      value: projects?.total || 0,
      icon: FolderKanban,
      color: 'bg-purple-500',
      href: '/admin/projects',
    },
    {
      name: 'Published Articles',
      value: articles?.data?.filter((a: any) => a.published).length || 0,
      icon: Eye,
      color: 'bg-green-500',
      href: '/admin/articles',
    },
    {
      name: 'Featured Projects',
      value: projects?.data?.filter((p: any) => p.featured).length || 0,
      icon: TrendingUp,
      color: 'bg-orange-500',
      href: '/admin/projects',
    },
  ];

  const StatCard = ({ stat }: { stat: typeof stats[0] }) => {
    const Icon = stat.icon;
    return (
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        whileHover={{ scale: 1.02, y: -4 }}
        className="card cursor-pointer group"
        onClick={() => navigate(stat.href)}
      >
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-1">
              {stat.name}
            </p>
            <p className="text-3xl font-bold text-gray-900 dark:text-white">
              {stat.value}
            </p>
          </div>
          <div className={`${stat.color} p-3 rounded-xl text-white`}>
            <Icon className="w-6 h-6" />
          </div>
        </div>
        <div className="mt-4 flex items-center text-sm text-primary-600 dark:text-primary-400 group-hover:gap-2 transition-all">
          <span>View all</span>
          <ArrowRight className="w-4 h-4 opacity-0 group-hover:opacity-100 transition-opacity" />
        </div>
      </motion.div>
    );
  };

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
          Dashboard
        </h1>
        <p className="text-gray-600 dark:text-gray-400">
          Welcome back! Here's what's happening with your portfolio.
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {stats.map((stat, index) => (
          <motion.div
            key={stat.name}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.1 }}
          >
            <StatCard stat={stat} />
          </motion.div>
        ))}
      </div>

      {/* Portfolio Info */}
      {portfolio && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          className="card"
        >
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
              Portfolio Information
            </h2>
            <button
              onClick={() => navigate('/admin/portfolio')}
              className="btn btn-secondary text-sm"
            >
              Edit Portfolio
            </button>
          </div>
          <div className="space-y-3">
            <div className="flex items-center gap-3">
              <User className="w-5 h-5 text-gray-400" />
              <div>
                <p className="font-medium text-gray-900 dark:text-white">{portfolio.name}</p>
                <p className="text-sm text-gray-600 dark:text-gray-400">{portfolio.title}</p>
              </div>
            </div>
            <p className="text-gray-600 dark:text-gray-400 text-sm">{portfolio.bio}</p>
          </div>
        </motion.div>
      )}

      {/* Recent Articles */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5 }}
          className="card"
        >
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
              Recent Articles
            </h2>
            <button
              onClick={() => navigate('/admin/articles')}
              className="text-sm text-primary-600 dark:text-primary-400 hover:underline"
            >
              View all
            </button>
          </div>
          <div className="space-y-4">
            {articles?.data?.slice(0, 3).map((article: any) => (
              <div
                key={article.id}
                className="p-4 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-primary-300 dark:hover:border-primary-700 transition-colors cursor-pointer"
                onClick={() => navigate('/admin/articles')}
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <h3 className="font-medium text-gray-900 dark:text-white mb-1">
                      {article.title}
                    </h3>
                    <p className="text-sm text-gray-600 dark:text-gray-400 line-clamp-2">
                      {article.excerpt}
                    </p>
                  </div>
                  {article.published && (
                    <span className="ml-4 px-2 py-1 text-xs font-medium bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-400 rounded">
                      Published
                    </span>
                  )}
                </div>
                <div className="mt-3 flex items-center gap-4 text-xs text-gray-500 dark:text-gray-400">
                  <div className="flex items-center gap-1">
                    <Calendar className="w-3 h-3" />
                    {new Date(article.created_at).toLocaleDateString()}
                  </div>
                </div>
              </div>
            ))}
            {(!articles?.data || articles.data.length === 0) && (
              <p className="text-center text-gray-500 dark:text-gray-400 py-8">
                No articles yet. Create your first article!
              </p>
            )}
          </div>
        </motion.div>

        {/* Recent Projects */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6 }}
          className="card"
        >
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
              Recent Projects
            </h2>
            <button
              onClick={() => navigate('/admin/projects')}
              className="text-sm text-primary-600 dark:text-primary-400 hover:underline"
            >
              View all
            </button>
          </div>
          <div className="space-y-4">
            {projects?.data?.slice(0, 3).map((project: any) => (
              <div
                key={project.id}
                className="p-4 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-primary-300 dark:hover:border-primary-700 transition-colors cursor-pointer"
                onClick={() => navigate('/admin/projects')}
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <h3 className="font-medium text-gray-900 dark:text-white mb-1">
                      {project.name}
                    </h3>
                    <p className="text-sm text-gray-600 dark:text-gray-400 line-clamp-2">
                      {project.description}
                    </p>
                  </div>
                  {project.featured && (
                    <span className="ml-4 px-2 py-1 text-xs font-medium bg-orange-100 dark:bg-orange-900/30 text-orange-800 dark:text-orange-400 rounded">
                      Featured
                    </span>
                  )}
                </div>
                <div className="mt-3 flex flex-wrap gap-2">
                  {project.technologies?.slice(0, 3).map((tech: string) => (
                    <span
                      key={tech}
                      className="px-2 py-1 text-xs bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded"
                    >
                      {tech}
                    </span>
                  ))}
                </div>
              </div>
            ))}
            {(!projects?.data || projects.data.length === 0) && (
              <p className="text-center text-gray-500 dark:text-gray-400 py-8">
                No projects yet. Add your first project!
              </p>
            )}
          </div>
        </motion.div>
      </div>
    </div>
  );
};

export default DashboardPage;
