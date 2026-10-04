import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import Link from 'next/link';
import { Star } from 'lucide-react';

interface Course {
  id: number;
  title: string;
  slug: string;
  short_description?: string;
  cover_url: string;
  difficulty?: string;
  rating?: number;
  reviews_count?: number;
}

export function CourseCard({ course, href }: { course: Course; href?: string }) {
  const linkHref = href || `/courses/${course.slug}`;
  const displayRating = course.rating ? Number(course.rating).toFixed(1) : '4.9';
  const reviewsCount = course.reviews_count ?? 16;

  return (
    <Link href={linkHref}>
      <Card className="overflow-hidden hover:shadow-xl transition-all duration-300 border-none bg-white dark:bg-slate-900 rounded-2xl h-full flex flex-col cursor-pointer group shadow-md">
        <div className="relative h-48 w-full overflow-hidden bg-gray-100 dark:bg-slate-800">
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            src={course.cover_url || 'https://images.unsplash.com/photo-1516321318423-f06f85e504b3?q=80&w=600&auto=format&fit=crop'}
            alt={course.title}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
          />
          {course.difficulty && (
            <div className="absolute top-3 left-3">
              <Badge variant="secondary" className="bg-white/90 dark:bg-slate-900/90 backdrop-blur-sm text-slate-800 dark:text-slate-200 border-none px-3 py-1 font-semibold text-xs">
                {course.difficulty === 'beginner'
                  ? 'Новичок'
                  : course.difficulty === 'intermediate'
                  ? 'Средний'
                  : course.difficulty === 'advanced'
                  ? 'Продвинутый'
                  : course.difficulty}
              </Badge>
            </div>
          )}
        </div>
        <CardHeader className="p-5 pb-2">
          <h3 className="font-bold text-lg line-clamp-2 text-slate-900 dark:text-white leading-tight">
            {course.title}
          </h3>
        </CardHeader>
        <CardContent className="p-5 pt-0 flex-grow flex flex-col justify-between space-y-3">
          <p className="text-xs text-slate-500 line-clamp-2">
            {course.short_description || 'Описание отсутствует'}
          </p>

          <div className="pt-2 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between text-xs font-semibold">
            <div className="flex items-center gap-1 text-amber-500">
              <Star size={13} className="fill-amber-400 text-amber-400" />
              <span>{displayRating}</span>
              <span className="text-slate-400 font-normal text-[11px]">
                ({reviewsCount} отзывов)
              </span>
            </div>
          </div>
        </CardContent>
      </Card>
    </Link>
  );
}
