import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { motion } from 'framer-motion';
import { Save, User, Mail, Github, Linkedin, Phone, Globe } from 'lucide-react';
import { portfolioService, Portfolio, UpdatePortfolioDto } from '../../services/api/portfolioService';

const PortfolioPage: React.FC = () => {
  const queryClient = useQueryClient();
  const [isEditing, setIsEditing] = useState(false);

  const { data: portfolio, isLoading } = useQuery({
    queryKey: ['portfolio'],
    queryFn: () => portfolioService.getPortfolio(),
  });

  const [error, setError] = useState<string>('');

  const updateMutation = useMutation({
    mutationFn: (portfolio: UpdatePortfolioDto) => portfolioService.updatePortfolio(portfolio),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['portfolio'] });
      setIsEditing(false);
      setError('');
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Failed to update portfolio');
    },
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
            Portfolio Settings
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Manage your portfolio information and social links
          </p>
        </div>
        {!isEditing && (
          <button
            onClick={() => setIsEditing(true)}
            className="btn btn-primary"
          >
            Edit Portfolio
          </button>
        )}
      </div>

      {/* Error Message */}
      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
          <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
          <button
            onClick={() => setError('')}
            className="mt-2 text-xs text-red-600 dark:text-red-400 hover:underline"
          >
            Dismiss
          </button>
        </div>
      )}

      {isEditing ? (
        <PortfolioForm
          portfolio={portfolio}
          onSubmit={(data) => {
            updateMutation.mutate(data);
          }}
          onCancel={() => setIsEditing(false)}
        />
      ) : (
        <PortfolioView portfolio={portfolio} />
      )}
    </div>
  );
};

const PortfolioView: React.FC<{ portfolio?: Portfolio }> = ({ portfolio }) => {
  if (!portfolio) {
    return (
      <div className="card text-center py-12">
        <p className="text-gray-500 dark:text-gray-400">No portfolio information available</p>
      </div>
    );
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="grid grid-cols-1 lg:grid-cols-2 gap-6"
    >
      {/* Basic Info */}
      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-6 flex items-center gap-2">
          <User className="w-5 h-5" />
          Basic Information
        </h2>
        <div className="space-y-4">
          <div>
            <label className="text-sm font-medium text-gray-500 dark:text-gray-400">Name</label>
            <p className="text-lg text-gray-900 dark:text-white mt-1">{portfolio.name}</p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-500 dark:text-gray-400">Title</label>
            <p className="text-lg text-gray-900 dark:text-white mt-1">{portfolio.title}</p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-500 dark:text-gray-400 flex items-center gap-2">
              <Mail className="w-4 h-4" />
              Email
            </label>
            <p className="text-lg text-gray-900 dark:text-white mt-1">{portfolio.email}</p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-500 dark:text-gray-400">Bio</label>
            <p className="text-gray-700 dark:text-gray-300 mt-1 leading-relaxed">{portfolio.bio}</p>
          </div>
        </div>
      </div>

      {/* Social Links */}
      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-6 flex items-center gap-2">
          <Globe className="w-5 h-5" />
          Social Links
        </h2>
        <div className="space-y-4">
          {portfolio.social_links?.github && (
            <a
              href={portfolio.social_links.github}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-3 p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-primary-300 dark:hover:border-primary-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-all"
            >
              <Github className="w-5 h-5 text-gray-600 dark:text-gray-400" />
              <span className="text-gray-900 dark:text-white">{portfolio.social_links.github}</span>
            </a>
          )}
          {portfolio.social_links?.linkedin && (
            <a
              href={portfolio.social_links.linkedin}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-3 p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-primary-300 dark:hover:border-primary-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-all"
            >
              <Linkedin className="w-5 h-5 text-gray-600 dark:text-gray-400" />
              <span className="text-gray-900 dark:text-white">{portfolio.social_links.linkedin}</span>
            </a>
          )}
          {portfolio.social_links?.phone && (
            <div className="flex items-center gap-3 p-3 rounded-lg border border-gray-200 dark:border-gray-700">
              <Phone className="w-5 h-5 text-gray-600 dark:text-gray-400" />
              <span className="text-gray-900 dark:text-white">{portfolio.social_links.phone}</span>
            </div>
          )}
        </div>
      </div>
    </motion.div>
  );
};

const PortfolioForm: React.FC<{
  portfolio?: Portfolio;
  onSubmit: (data: UpdatePortfolioDto) => void;
  onCancel: () => void;
}> = ({ portfolio, onSubmit, onCancel }) => {
  const [name, setName] = useState(portfolio?.name || '');
  const [title, setTitle] = useState(portfolio?.title || '');
  const [email, setEmail] = useState(portfolio?.email || '');
  const [bio, setBio] = useState(portfolio?.bio || '');
  const [github, setGithub] = useState(portfolio?.social_links?.github || '');
  const [linkedin, setLinkedin] = useState(portfolio?.social_links?.linkedin || '');
  const [phone, setPhone] = useState(portfolio?.social_links?.phone || '');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      name,
      title,
      email,
      bio,
      social_links: {
        github,
        linkedin,
        phone,
      },
    });
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="card"
    >
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="label">Full Name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              className="input"
              placeholder="John Doe"
            />
          </div>
          <div>
            <label className="label">Professional Title</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              className="input"
              placeholder="Fullstack Developer"
            />
          </div>
        </div>

        <div>
          <label className="label flex items-center gap-2">
            <Mail className="w-4 h-4" />
            Email Address
          </label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="input"
            placeholder="your.email@example.com"
          />
        </div>

        <div>
          <label className="label">Bio</label>
          <textarea
            value={bio}
            onChange={(e) => setBio(e.target.value)}
            required
            className="input"
            rows={6}
            placeholder="Write a brief description about yourself..."
          />
        </div>

        <div className="border-t border-gray-200 dark:border-gray-700 pt-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Social Links</h3>
          <div className="space-y-4">
            <div>
              <label className="label flex items-center gap-2">
                <Github className="w-4 h-4" />
                GitHub URL
              </label>
              <input
                type="url"
                value={github}
                onChange={(e) => setGithub(e.target.value)}
                className="input"
                placeholder="https://github.com/username"
              />
            </div>
            <div>
              <label className="label flex items-center gap-2">
                <Linkedin className="w-4 h-4" />
                LinkedIn URL
              </label>
              <input
                type="url"
                value={linkedin}
                onChange={(e) => setLinkedin(e.target.value)}
                className="input"
                placeholder="https://linkedin.com/in/username"
              />
            </div>
            <div>
              <label className="label flex items-center gap-2">
                <Phone className="w-4 h-4" />
                Phone Number
              </label>
              <input
                type="tel"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                className="input"
                placeholder="+1234567890"
              />
            </div>
          </div>
        </div>

        <div className="flex items-center justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
          <button
            type="button"
            onClick={onCancel}
            className="btn btn-secondary"
          >
            Cancel
          </button>
          <button
            type="submit"
            className="btn btn-primary flex items-center gap-2"
          >
            <Save className="w-4 h-4" />
            Save Changes
          </button>
        </div>
      </form>
    </motion.div>
  );
};

export default PortfolioPage;
