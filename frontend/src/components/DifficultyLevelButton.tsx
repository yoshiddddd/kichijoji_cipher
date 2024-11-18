import React from "react";

type DifficultyLevelButtonProps = {
  label: string;
  roomLevel: number;
  changeRoomLevel: (level: number) => void;
};

const DifficultyLevelButton = ({
  label,
  roomLevel,
  changeRoomLevel,
}: DifficultyLevelButtonProps) => {
  const setRoomLevel = (label: string) => {
    switch (label) {
      case "初級":
        return 1;
      case "中級":
        return 2;
      case "上級":
        return 3;
      default:
        return 0;
    }
  };
  return (
    <button
      onClick={() => changeRoomLevel(setRoomLevel(label))}
      style={{
        backgroundColor:
          roomLevel === setRoomLevel(label) ? "lightblue" : "white",
      }}
    >
      {label}
    </button>
  );
};

export default DifficultyLevelButton;
