#!/bin/bash
cd "c:\Git Repos\Random Repos\roadmap"
echo "Running ESLint on roadmap files..."
npx eslint public/AsyncPages.tsx public/pages/Administration/pages/ManageRoadmap.page.tsx public/pages/Roadmap/Roadmap.page.tsx public/pages/Roadmap/components/AssignToRoadmapModal.tsx public/pages/Roadmap/components/RoadmapColumn.tsx public/pages/Roadmap/components/RoadmapPostCard.tsx public/pages/ShowPost/ShowPost.page.tsx public/services/roadmap.ts --fix
echo "ESLint completed."
