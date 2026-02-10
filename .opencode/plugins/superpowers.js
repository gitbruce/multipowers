import fs from 'fs';
import path from 'path';
import { findSkillsInDir, resolveSkillPath } from '../../lib/skills-core.js';

function listInstalledSkills(superpowersDir, personalDir) {
    const superpowerSkills = findSkillsInDir(superpowersDir, 'superpowers', 4);
    const personalSkills = findSkillsInDir(personalDir, 'personal', 4);

    return [...personalSkills, ...superpowerSkills].map((item) => ({
        name: item.name,
        description: item.description,
        sourceType: item.sourceType,
        skillFile: item.skillFile,
    }));
}

function buildRuntime(superpowersRoot, personalRoot) {
    const superpowersSkillsDir = path.join(superpowersRoot, 'skills');
    const personalSkillsDir = path.join(personalRoot, 'skills');

    return {
        listSkills() {
            return listInstalledSkills(superpowersSkillsDir, personalSkillsDir);
        },
        resolveSkill(name) {
            return resolveSkillPath(name, superpowersSkillsDir, personalSkillsDir);
        },
        health() {
            const skills = listInstalledSkills(superpowersSkillsDir, personalSkillsDir);
            return {
                plugin: 'superpowers',
                superpowersRoot,
                personalRoot,
                skillCount: skills.length,
            };
        },
    };
}

export default {
    name: 'superpowers',
    description: 'Loads superpowers and personal skills with deterministic precedence.',
    setup(context = {}) {
        const home = process.env.HOME || process.cwd();
        const superpowersRoot = context.superpowersRoot || path.join(home, '.config/opencode/superpowers');
        const personalRoot = context.personalRoot || path.join(home, '.config/opencode');

        if (!fs.existsSync(superpowersRoot)) {
            throw new Error(`superpowers root not found: ${superpowersRoot}`);
        }

        return buildRuntime(superpowersRoot, personalRoot);
    },
};
