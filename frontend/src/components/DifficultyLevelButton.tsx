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
          roomLevel === setRoomLevel(label)
            ? label === "初級"
              ? "#90EE90"
              : label === "中級"
              ? "#FFD700"
              : "#FF6B6B"
            : "white",
        padding: "10px 20px",
        margin: "0 10px",
        border:
          roomLevel === setRoomLevel(label)
            ? label === "初級"
              ? "2px solid #32CD32"
              : label === "中級"
              ? "2px solid #DAA520"
              : "2px solid #DC143C"
            : "1px solid #ccc",
        borderRadius: "5px",
        cursor: "pointer",
        fontWeight: roomLevel === setRoomLevel(label) ? "bold" : "normal",
        color:
          roomLevel === setRoomLevel(label)
            ? label === "初級"
              ? "#006400"
              : label === "中級"
              ? "#8B4513"
              : "#8B0000"
            : "#333",
        transition: "all 0.3s ease",
      }}
    >
      {label}
    </button>
  );
};

export default DifficultyLevelButton;
