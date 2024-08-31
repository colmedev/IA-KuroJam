import { useAuth } from "@clerk/clerk-react";
import { useCallback, useEffect, useState } from "react";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

type Career = {
  id: number;
  title: string;
  description: string;
  personalityDescription: string;
  education: string;
  averageSalary: string;
  lowerSalary: string;
  highestSalary: string;
  tasks: string[] | null;
  careerTasks: Task[];
  knowledge: KnowledgeCategory[] | null;
  abilities: AbilityCategory[] | null;
  skills: SkillCategory[] | null;
  technology: TechnologyCategory[] | null;
  personality: Personality;
};

type Task = {
  id: number;
  careerId: number;
  taskDescription: string;
};

type KnowledgeCategory = {
  id: number;
  careerId: number;
  name: string;
  areas: string[];
};

type AbilityCategory = {
  id: number;
  careerId: number;
  name: string;
  areas: string[];
};

type SkillCategory = {
  id: number;
  careerId: number;
  name: string;
  areas: string[];
};

type TechnologyCategory = {
  id: number;
  careerId: number;
  name: string;
  areas: string[];
};

type Personality = {
  description: string;
  attributes: string[] | null;
};

const Results = () => {

  const { getToken } = useAuth();
  const [careers, setCareers] = useState<Career[]>([]);

  const fetchResults = useCallback(async () => {
    try {
      const token = await getToken();

      const response = await fetch(`${BACKEND_URL}/results`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error('Error loading results');
      }

      const responseData = await response.json();

      console.log(responseData.careers);
      setCareers(responseData.careers);
    } catch (error) {
      alert(error);
    }
  }, [getToken]);

  useEffect(() => {
    fetchResults();
  }, [fetchResults])

  return (
    <div className="results-page">
      {careers.map((career) => (
        <div key={career.id} className="career-card">
          <h2>{career.title}</h2>
          <p><strong>Descripción:</strong> {career.description}</p>
          <p><strong>Descripción de la Personalidad:</strong> {career.personalityDescription}</p>
          <p><strong>Educación:</strong> {career.education}</p>
          <p><strong>Salario Promedio:</strong> {career.averageSalary}</p>
          <p><strong>Rango Salarial:</strong> {career.lowerSalary} - {career.highestSalary}</p>

          <h3>Tareas</h3>
          <ul>
            {career.careerTasks?.map((task) => (
              <li key={task.id}>{task.taskDescription}</li>
            ))}
          </ul>

          <h3>Conocimientos</h3>
          {career.knowledge ? (
            <ul>
              {career.knowledge.map((knowledge) => (
                <li key={knowledge.id}>
                  <strong>{knowledge.name}:</strong> {knowledge.areas.join(", ")}
                </li>
              ))}
            </ul>
          ) : (
            <p>No hay datos de conocimientos disponibles</p>
          )}

          <h3>Habilidades</h3>
          {career.abilities ? (
            <ul>
              {career.abilities.map((ability) => (
                <li key={ability.id}>
                  <strong>{ability.name}:</strong> {ability.areas.join(", ")}
                </li>
              ))}
            </ul>
          ) : (
            <p>No hay datos de habilidades disponibles</p>
          )}

          <h3>Categorías de Habilidades</h3>
          {career.skills ? (
            <ul>
              {career.skills.map((skillCategory) => (
                <li key={skillCategory.id}>
                  <strong>{skillCategory.name}:</strong> {skillCategory.areas.join(", ")}
                </li>
              ))}
            </ul>
          ) : (
            <p>No hay datos de categorías de habilidades disponibles</p>
          )}

          <h3>Categorías de Tecnología</h3>
          {career.technology ? (
            <ul>
              {career.technology.map((technologyCategory) => (
                <li key={technologyCategory.id}>
                  <strong>{technologyCategory.name}:</strong> {technologyCategory.areas.join(", ")}
                </li>
              ))}
            </ul>
          ) : (
            <p>No hay datos de categorías de tecnología disponibles</p>
          )}

          <h3>Atributos de Personalidad</h3>
          {career.personality.attributes ? (
            <ul>
              {career.personality.attributes.map((attribute, index) => (
                <li key={index}>
                  <strong>{attribute}</strong>
                </li>
              ))}
            </ul>
          ) : (
            <p>No hay atributos de personalidad disponibles</p>
          )}
        </div>
      ))}
    </div>
  );
};

export default Results;
