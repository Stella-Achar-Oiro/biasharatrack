import { BookOpen, PlayCircle } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';

export default function TutorialPanel() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-[#011627]">{t('tutorialPanel.title')}</h2>
        <BookOpen className="h-6 w-6 text-[#2EC4B6]" />
      </div>
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <span className="text-sm text-gray-500">{t('tutorialPanel.newVideos')}</span>
          <span className="text-lg font-medium">8</span>
        </div>
        <div className="flex items-center text-[#2EC4B6]">
          <PlayCircle className="h-5 w-5 mr-2" />
          <span className="text-sm">{t('tutorialPanel.latestVideo')}</span>
        </div>
        <button 
          onClick={() => navigate('/dashboard/tutorials')}
          className="w-full bg-[#2EC4B6] text-white py-2 rounded-lg hover:bg-[#28b0a3] transition-colors"
        >
          {t('tutorialPanel.button')}
        </button>
      </div>
    </div>
  );
} 