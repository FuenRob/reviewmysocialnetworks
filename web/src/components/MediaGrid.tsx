import React, { useState } from 'react';
import type { MediaAnalysisItem } from '../types/instagram';
import { Heart, MessageCircle, ExternalLink, Video, Layers, Image as ImageIcon, Flame, Bookmark, Eye, Share2 } from 'lucide-react';

interface Props {
  media: MediaAnalysisItem[];
  platform: 'instagram' | 'tiktok';
}

export const MediaGrid: React.FC<Props> = ({ media, platform }) => {
  const [activeFilter, setActiveFilter] = useState<'ALL' | 'VIDEO' | 'CAROUSEL_ALBUM' | 'IMAGE'>('ALL');

  const filteredMedia = media.filter((item) => {
    if (activeFilter === 'ALL') return true;
    return item.media_type === activeFilter;
  });

  const getFormatBadge = (type: string) => {
    switch (type) {
      case 'VIDEO':
        return {
          icon: Video,
          label: platform === 'tiktok' ? 'Vídeo TikTok' : 'Reel / Vídeo',
          color: 'bg-pink-500/20 text-pink-300 border-pink-500/30',
        };
      case 'CAROUSEL_ALBUM':
        return {
          icon: Layers,
          label: 'Carrusel',
          color: 'bg-purple-500/20 text-purple-300 border-purple-500/30',
        };
      case 'IMAGE':
      default:
        return {
          icon: ImageIcon,
          label: 'Foto',
          color: 'bg-blue-500/20 text-blue-300 border-blue-500/30',
        };
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h3 className="text-lg font-bold text-white">
            Publicaciones Analizadas ({media.length})
          </h3>
          <p className="text-xs text-slate-400">
            Detalle individual de engagement, métricas e interacciones
          </p>
        </div>

        {platform === 'instagram' && <div className="flex items-center gap-1.5 p-1 bg-slate-900 border border-slate-800 rounded-2xl w-fit">
          {[
            { key: 'ALL', label: 'Todos' },
            { key: 'CAROUSEL_ALBUM', label: 'Carruseles' },
            { key: 'VIDEO', label: 'Reels' },
            { key: 'IMAGE', label: 'Fotos' },
          ].map((tab) => (
            <button
              key={tab.key}
              onClick={() => setActiveFilter(tab.key as any)}
              className={`px-3 py-1.5 rounded-xl text-xs font-semibold transition-all ${
                activeFilter === tab.key
                  ? 'bg-indigo-600 text-white shadow-md'
                  : 'text-slate-400 hover:text-slate-200'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>}
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredMedia.map((item) => {
          const badge = getFormatBadge(item.media_type);
          const BadgeIcon = badge.icon;
          const dateStr = new Date(item.timestamp).toLocaleDateString('es-ES', {
            day: 'numeric',
            month: 'short',
            year: 'numeric',
          });
          const displayedEngagement = platform === 'tiktok' ? item.view_engagement_rate || 0 : item.engagement_rate;

          return (
            <div
              key={item.id}
              className="group bg-slate-900/80 border border-slate-800 rounded-3xl overflow-hidden backdrop-blur-xl flex flex-col justify-between hover:border-slate-700 transition-all hover:shadow-xl"
            >
              <div className="relative aspect-square w-full bg-slate-950 overflow-hidden">
                <img
                  src={item.thumbnail_url || item.media_url || 'https://images.unsplash.com/photo-1517694712202-14dd9538aa97?w=600&auto=format&fit=crop&q=80'}
                  alt={item.caption || `Publicación de ${platform === 'tiktok' ? 'TikTok' : 'Instagram'}`}
                  className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                />

                <div className="absolute top-3 left-3 right-3 flex items-center justify-between pointer-events-none">
                  <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-bold border backdrop-blur-md ${badge.color}`}>
                    <BadgeIcon className="w-3 h-3" />
                    {badge.label}
                  </span>

                  {item.is_top_performer && (
                    <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-bold bg-amber-500 text-slate-950 shadow-lg shadow-amber-500/30">
                      <Flame className="w-3 h-3 fill-slate-950" />
                      Top Post
                    </span>
                  )}
                </div>

                <div className="absolute bottom-0 inset-x-0 bg-gradient-to-t from-slate-950 via-slate-950/80 to-transparent p-3 pt-6 flex items-center justify-between text-xs">
                  <div className="flex items-center gap-3 text-white font-semibold">
                    <span className="flex items-center gap-1">
                      <Heart className="w-3.5 h-3.5 text-rose-400 fill-rose-400/20" />
                      {item.like_count.toLocaleString()}
                    </span>
                    {platform === 'tiktok' && <span className="flex items-center gap-1"><Eye className="w-3.5 h-3.5 text-cyan-400" />{(item.view_count || 0).toLocaleString()}</span>}
                    {platform === 'tiktok' && <span className="flex items-center gap-1"><Share2 className="w-3.5 h-3.5 text-purple-400" />{(item.share_count || 0).toLocaleString()}</span>}
                    <span className="flex items-center gap-1">
                      <MessageCircle className="w-3.5 h-3.5 text-blue-400 fill-blue-400/20" />
                      {item.comments_count.toLocaleString()}
                    </span>
                    {item.insights?.saved ? (
                      <span className="flex items-center gap-1 text-purple-400">
                        <Bookmark className="w-3.5 h-3.5 fill-purple-400/20" />
                        {item.insights.saved}
                      </span>
                    ) : null}
                  </div>

                  <span className="px-2 py-0.5 rounded-md bg-emerald-500/20 border border-emerald-500/30 text-emerald-400 font-bold text-xs">
                    {displayedEngagement}% eng
                  </span>
                </div>
              </div>

              <div className="p-4 space-y-3 flex-1 flex flex-col justify-between">
                <div>
                  <span className="text-[10px] text-slate-500 block mb-1">
                    Publicado el {dateStr}{item.duration_seconds ? ` · ${item.duration_seconds}s` : ''}
                  </span>
                  <p className="text-xs text-slate-300 line-clamp-3 leading-relaxed">
                    {item.caption || '(Sin texto de pie de foto)'}
                  </p>
                </div>

                {item.permalink && (
                  <div className="pt-2 border-t border-slate-800/80 flex justify-end">
                    <a
                      href={item.permalink}
                      target="_blank"
                      rel="noreferrer"
                      className="inline-flex items-center gap-1 text-[11px] font-semibold text-indigo-400 hover:text-indigo-300 transition-colors"
                    >
                      <span>Ver en {platform === 'tiktok' ? 'TikTok' : 'Instagram'}</span>
                      <ExternalLink className="w-3 h-3" />
                    </a>
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
