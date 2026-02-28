export interface ActionDescription {
  command: string;
  description: string;
}

export interface LessonStageSummary {
  index: number;
  id: string;
  title: string;
  objective?: string;
  goal?: string;
  passConditions?: string[];
  prerequisites?: string[];
  allowedCommands?: string[];
  actionDescriptions?: ActionDescription[];
  expectedVisualCues?: string[];
  limits?: {
    maxSteps?: number;
    maxPolicyChanges?: number;
    maxConfigChanges?: number;
  };
  attempts?: number;
  completed?: boolean;
  unlocked?: boolean;
}

export interface LessonSummary {
  id: string;
  title: string;
  module: string;
  sectionId?: string;
  sectionTitle?: string;
  difficulty?: string;
  estimatedMinutes?: number;
  chapterRefs?: string[];
  stages: LessonStageSummary[];
}

export interface CurriculumSection {
  id: string;
  title: string;
  subtitle?: string;
  order: number;
  comingSoon: boolean;
  lessons?: LessonSummary[];
  completedStages?: number;
  totalStages?: number;
}

export interface LessonLearnStage {
  index: number;
  id: string;
  title: string;
  coreIdea?: string;
  mechanismSteps?: string[];
  workedExample?: string;
  commonMistakes?: string[];
  preChallengeChecklist?: string[];
  objective?: string;
  goal?: string;
  prerequisites?: string[];
  expectedVisualCues?: string[];
}

export interface LessonLearnSummary {
  id: string;
  title: string;
  module: string;
  sectionId?: string;
  sectionTitle?: string;
  difficulty?: string;
  estimatedMinutes?: number;
  chapterRefs?: string[];
  stages: LessonLearnStage[];
}
