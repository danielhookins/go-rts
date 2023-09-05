package main

import (
    raylib "github.com/gen2brain/raylib-go/raylib"
    "math"
)

const (
    screenWidth  = 800
    screenHeight = 600
    unitSize     = 20
)

type Unit struct {
    Position raylib.Vector2
    Color    raylib.Color
    Selected bool
    Speed    float32
}

func avoidOverlap(units []Unit) {
    for i := range units {
        for j := range units {
            if i != j {
                if raylib.CheckCollisionRecs(
                    raylib.NewRectangle(units[i].Position.X, units[i].Position.Y, unitSize, unitSize),
                    raylib.NewRectangle(units[j].Position.X, units[j].Position.Y, unitSize, unitSize),
                ) {
                    dir := raylib.Vector2Subtract(units[i].Position, units[j].Position)
                    length := raylib.Vector2Length(dir)
                    if length == 0 {
                        length = 1
                    }
                    dir = raylib.Vector2Scale(dir, (unitSize-length)/length)
                    units[i].Position = raylib.Vector2Add(units[i].Position, dir)
                }
            }
        }
    }
}

func main() {
    raylib.InitWindow(screenWidth, screenHeight, "RTS Game")
    raylib.SetTargetFPS(60)

    playerUnits := []Unit{
        {raylib.NewVector2(100, 100), raylib.Blue, false, 2.0},
        {raylib.NewVector2(200, 200), raylib.Blue, false, 2.0},
    }

    enemyUnits := []Unit{
        {raylib.NewVector2(600, 100), raylib.Red, false, 0},
        {raylib.NewVector2(700, 200), raylib.Red, false, 0},
    }

    selectRectangle := raylib.NewRectangle(0, 0, 0, 0)
    selecting := false
    target := raylib.NewVector2(0, 0)
    hasTarget := false

    for !raylib.WindowShouldClose() {
        raylib.BeginDrawing()
        raylib.ClearBackground(raylib.RayWhite)

        // Draw units
        for _, unit := range playerUnits {
            raylib.DrawRectangleV(unit.Position, raylib.NewVector2(unitSize, unitSize), unit.Color)
            if unit.Selected {
                raylib.DrawRectangleLinesEx(raylib.NewRectangle(unit.Position.X, unit.Position.Y, unitSize, unitSize), 2, raylib.DarkBlue)
            }
        }

        for _, unit := range enemyUnits {
            raylib.DrawRectangleV(unit.Position, raylib.NewVector2(unitSize, unitSize), unit.Color)
        }

        if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
            if !selecting {
                selectRectangle.X = float32(raylib.GetMouseX())
                selectRectangle.Y = float32(raylib.GetMouseY())
                selecting = true
            } else {
                endX := float32(raylib.GetMouseX())
                endY := float32(raylib.GetMouseY())
                
                if endX < selectRectangle.X {
                    selectRectangle.Width = selectRectangle.X - endX
                    selectRectangle.X = endX
                } else {
                    selectRectangle.Width = endX - selectRectangle.X
                }
                
                if endY < selectRectangle.Y {
                    selectRectangle.Height = selectRectangle.Y - endY
                    selectRectangle.Y = endY
                } else {
                    selectRectangle.Height = endY - selectRectangle.Y
                }
                
                raylib.DrawRectangleLinesEx(selectRectangle, 1, raylib.Gray)
            }
        }        

        if raylib.IsMouseButtonReleased(raylib.MouseLeftButton) {
            if selecting {
                hasTarget = false // Reset the target
                for i := range playerUnits {
                    playerUnits[i].Selected = raylib.CheckCollisionRecs(selectRectangle, raylib.NewRectangle(playerUnits[i].Position.X, playerUnits[i].Position.Y, unitSize, unitSize))
                }
                selecting = false
            }
        }

        // Set target for selected units
        if raylib.IsMouseButtonPressed(raylib.MouseRightButton) {
            target = raylib.NewVector2(float32(raylib.GetMouseX()), float32(raylib.GetMouseY()))
            hasTarget = true
        }

        // Move selected units to mouse position
        if hasTarget {
            for i := range playerUnits {
                if playerUnits[i].Selected {
                    distanceToTarget := raylib.Vector2Distance(playerUnits[i].Position, target)
                    if distanceToTarget > playerUnits[i].Speed {
                        angle := math.Atan2(float64(target.Y-playerUnits[i].Position.Y), float64(target.X-playerUnits[i].Position.X))
                        playerUnits[i].Position.X += float32(math.Cos(angle) * float64(playerUnits[i].Speed))
                        playerUnits[i].Position.Y += float32(math.Sin(angle) * float64(playerUnits[i].Speed))
                    }
                }
            }
            avoidOverlap(playerUnits)
        }

        raylib.EndDrawing()
    }

    raylib.CloseWindow()
}

