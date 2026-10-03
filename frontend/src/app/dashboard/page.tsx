'use client';

import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { CourseCard } from '@/components/CourseCard';
import { BookOpen } from 'lucide-react';

export default function StudentDashboard() {
  const [courses, setCourses] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchMyCourses = async () => {
      try {
        const { data } = await api.get('/courses/my');
        setCourses(data.data?.courses || []);
      } catch (err) {
        console.error('Failed to load user courses', err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchMyCourses();
  }, []);

  if (isLoading) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-extrabold text-gray-900">Мое обучение</h1>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {[1, 2, 3].map(i => <div key={i} className="h-72 bg-gray-200 animate-pulse rounded-2xl"></div>)}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-extrabold text-gray-900">Мое обучение</h1>
        <p className="text-gray-500 mt-2">Продолжайте изучать начатые курсы</p>
      </div>

      {courses.length === 0 ? (
        <div className="flex flex-col items-center justify-center p-12 bg-white rounded-3xl border border-dashed border-gray-200">
          <BookOpen size={48} className="text-gray-300 mb-4" />
          <h2 className="text-xl font-bold text-gray-700">У вас пока нет курсов</h2>
          <p className="text-gray-500 mt-2 mb-6">Перейдите в каталог, чтобы найти что-то интересное.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {courses.map((course: any) => (
            <CourseCard 
              key={course.id} 
              course={course} 
              href={`/dashboard/courses/${course.id}`} // Направляем внутрь дашборда курса
            />
          ))}
        </div>
      )}
    </div>
  );
}
