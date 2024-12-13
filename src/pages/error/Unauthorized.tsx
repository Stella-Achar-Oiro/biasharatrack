import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { useTranslation } from 'react-i18next';

export default function Unauthorized() {
  const navigate = useNavigate();
  const { t } = useTranslation();

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <motion.h1 
          className="text-9xl font-bold text-[#011627]"
          initial={{ opacity: 0, y: -50 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          401
        </motion.h1>
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.2, duration: 0.5 }}
        >
          <p className="text-2xl font-semibold text-gray-600 mt-4">{t('errors.unauthorized.title')}</p>
          <p className="text-gray-500 mt-2">{t('errors.unauthorized.message')}</p>
        </motion.div>
        <motion.div 
          className="mt-6"
          initial={{ opacity: 0, y: 50 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4, duration: 0.5 }}
        >
          <motion.button
            onClick={() => navigate('/login')}
            className="mr-4 px-4 py-2 bg-[#2EC4B6] text-white rounded-lg hover:bg-[#28b0a3] transition-colors"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            {t('errors.unauthorized.buttons.login')}
          </motion.button>
          <motion.button
            onClick={() => navigate('/')}
            className="px-4 py-2 border border-[#2EC4B6] text-[#2EC4B6] rounded-lg hover:bg-[#2EC4B6] hover:text-white transition-colors"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            {t('errors.unauthorized.buttons.home')}
          </motion.button>
        </motion.div>
      </div>
    </div>
  );
} 