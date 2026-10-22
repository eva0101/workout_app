export type Workout = {
  id: string
  name: string
  date: string
  duration: string
  exercises: number
  tone: 'sage' | 'sand' | 'blue'
}

export const recentWorkouts: Workout[] = [
  {
    id: '1',
    name: 'Верх тела',
    date: 'Сегодня, 09:24',
    duration: '52 мин',
    exercises: 6,
    tone: 'sage',
  },
  {
    id: '2',
    name: 'Ноги и пресс',
    date: 'Воскресенье, 11:10',
    duration: '1 ч 08 мин',
    exercises: 7,
    tone: 'sand',
  },
  {
    id: '3',
    name: 'Спина и бицепс',
    date: 'Пятница, 18:42',
    duration: '46 мин',
    exercises: 5,
    tone: 'blue',
  },
]

export const nextWorkout = {
  name: 'Ноги и плечи',
  focus: 'СИЛА · ТРЕНИРОВКА B',
  duration: 'около 55 мин',
  exercises: [
    'Приседания со штангой',
    'Румынская тяга',
    'Жим гантелей сидя',
    'Подъём на носки',
    'Жим ногами',
    'Махи гантелями в стороны',
  ],
}
