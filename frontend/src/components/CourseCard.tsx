import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import Link from 'next/link';

interface Course {
  id: number;
  title: string;
  slug: string;
  short_description?: string;
  cover_url: string;
  difficulty?: string;
}

export function CourseCard({ course, href }: { course: Course, href?: string }) {
  const linkHref = href || `/courses/${course.slug}`;
  
  return (
    <Link href={linkHref}>
      <Card className="overflow-hidden hover:shadow-xl transition-all duration-300 border-none bg-white rounded-2xl h-full flex flex-col cursor-pointer group shadow-md">
        <div className="relative h-48 w-full overflow-hidden bg-gray-100">
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img 
            src={course.cover_url || 'https://images.unsplash.com/photo-1516321318423-f06f85e504b3?q=80&w=600&auto=format&fit=crop'} 
            alt={course.title} 
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
          />
          {course.difficulty && (
            <div className="absolute top-3 left-3">
              <Badge variant="secondary" className="bg-white/90 backdrop-blur-sm text-black border-none px-3 py-1 font-semibold">
                {course.difficulty === 'beginner' ? 'Новичок' : 
                 course.difficulty === 'intermediate' ? 'Средний' : 
                 course.difficulty === 'advanced' ? 'Продвинутый' : course.difficulty}
              </Badge>
            </div>
          )}
        </div>
        <CardHeader className="p-5 pb-3">
          <h3 className="font-bold text-xl line-clamp-2 text-gray-900 leading-tight">{course.title}</h3>
        </CardHeader>
        <CardContent className="p-5 pt-0 flex-grow">
          <p className="text-sm text-muted-foreground line-clamp-3">
            {course.short_description || 'Описание отсутствует'}
          </p>
        </CardContent>
      </Card>
    </Link>
  );
}
