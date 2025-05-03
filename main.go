package main

import (
	"DBD_QTE_go/http_api"
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	PointerSpeed    = 120
	NormalQTEAngle  = 30
	PerfectQTEAngle = 10

	PerfectQTEProgressBonus  = 3
	FailureProgressRetreat   = 1
	FailureProgressPauseTime = 300 * time.Millisecond
	QTETipSoundDuration      = 500 * time.Millisecond
	frameRate                = 120
	rotationSpeedPerFrame    = PointerSpeed / float32(frameRate)
	LongestWaitingTime       = time.Second * 15
	ShortestWaitingTime      = time.Second * 5
	//clientId                 = "6c13674b-9960-4ae6-8165-a53a4752c5e0"
)

type GameInfoResponse struct {
	Status         int    `json:"status"`
	Code           string `json:"code"`
	StrengthConfig struct {
		Strength       int `json:"strength"`
		RandomStrength int `json:"randomStrength"`
	} `json:"strengthConfig"`
}

var gameInfo GameInfoResponse
var clientId string

func init() {
	if LongestWaitingTime < ShortestWaitingTime {
		_ = fmt.Sprintf("Error: LongestWaitingTime (%v) is less than ShortestWaitingTime (%v)", LongestWaitingTime, ShortestWaitingTime)
		_ = fmt.Sprintf("错误: 最长等待时间 (%v) 小于最短等待时间 (%v)", LongestWaitingTime, ShortestWaitingTime)
		log.Fatalf("Error: LongestWaitingTime (%v) is less than ShortestWaitingTime (%v)\n错误: 最长等待时间 (%v) 小于最短等待时间 (%v)", LongestWaitingTime, ShortestWaitingTime, LongestWaitingTime, ShortestWaitingTime)
	}
	// 获取用户输入的 clientId
	fmt.Print("请输入 clientId: ")
	_, err := fmt.Scanln(&clientId)
	if err != nil {
		log.Fatalf("无法读取 clientId: %v", err)
	}

	gameInfo, err := http_api.GetGameInfo(clientId) // 从 API 获取游戏信息
	if err != nil {
		log.Fatalf("获取游戏信息失败: %v", err)
	}
	// 输出获取的基础强度
	fmt.Printf("获取到的基础强度: %d\n", gameInfo.StrengthConfig.Strength)
}
func DrawArc(center rl.Vector2, radius float32, startAngle, endAngle float32, lineThick float32, col rl.Color) {
	const segments = 200
	step := (endAngle - startAngle) / float32(segments) // Calculate the step size between points
	angle := startAngle

	// Draw each line segment along the arc
	for i := 0; i < segments; i++ {
		// Calculate the starting point of the segment
		startX := center.X + radius*float32(math.Cos(float64(angle)*math.Pi/180))
		startY := center.Y + radius*float32(math.Sin(float64(angle)*math.Pi/180))

		// Calculate the ending point of the segment
		angle += step
		endX := center.X + radius*float32(math.Cos(float64(angle)*math.Pi/180))
		endY := center.Y + radius*float32(math.Sin(float64(angle)*math.Pi/180))

		// Draw the line segment (arc)
		rl.DrawLine(int32(startX), int32(startY), int32(endX), int32(endY), col)
	}
}

func DrawArcCommonQTE(center rl.Vector2, radius float32, startAngle, endAngle float32, lineThick float32, col rl.Color) {
	// Draw the common arc at a fixed angle, independent of pointer rotation
	for i := -2; i < 2; i++ {
		DrawArc(center, radius+float32(i), startAngle, endAngle, lineThick, col)
	}
}

func DrawArcPerfectQTE(center rl.Vector2, radius float32, startAngle, endAngle float32, lineThick float32, col rl.Color) {
	// Draw the perfect arc at a fixed angle, independent of pointer rotation
	for i := -5; i < 5; i++ {
		DrawArc(center, radius+float32(i), startAngle, endAngle, lineThick, col)
	}
}
func DrawPointerRotation(center rl.Vector2, radius float32, startAngle float32, direction string) float32 {
	// Adjust angle based on rotation direction
	angle := startAngle
	if direction == "clockwise" {
		angle += rotationSpeedPerFrame
	} else if direction == "counterclockwise" {
		angle -= rotationSpeedPerFrame
	}

	// Ensure angle stays within 0-360 degrees
	if angle >= 360 {
		angle -= 360
	} else if angle < 0 {
		angle += 360
	}

	PointerX := center.X + radius*float32(math.Cos(float64(angle)*math.Pi/180))
	PointerY := center.Y + radius*float32(math.Sin(float64(angle)*math.Pi/180))

	// Draw the rotating pointer line
	rl.DrawLine(int32(center.X), int32(center.Y), int32(PointerX), int32(PointerY), rl.Blue)

	return angle
}
func triggerFunction() {
	// 当指针旋转 360 度时触发的函数
	fmt.Println("Pointer has rotated 360 degrees, triggering function!")

	gameInfo, err := http_api.GetGameInfo(clientId) // 从 API 获取游戏信息
	if err != nil {
		log.Fatalf("获取游戏信息失败: %v", err)
	}
	// 增加基础强度
	gameInfo.StrengthConfig.Strength += FailureProgressRetreat
	fmt.Printf("基础强度加: %d\n", FailureProgressRetreat)

	// 创建 SetStrengthConfigRequest 请求体
	config := http_api.SetStrengthConfigRequest{
		Strength: struct {
			Add *int `json:"add,omitempty"`
			Sub *int `json:"sub,omitempty"`
			Set *int `json:"set,omitempty"`
		}{
			Set: &gameInfo.StrengthConfig.Strength, // 传递增加的强度
		},
	}

	// 调用 http_api.SetStrength 更新基础强度
	err = http_api.SetStrength(clientId, config)
	if err != nil {
		log.Printf("更新基础强度失败: %v", err)
		return
	}
	fmt.Printf("成功更新基础强度: %d\n", gameInfo.StrengthConfig.Strength)
}
func perfectQTE() {
	gameInfo, err := http_api.GetGameInfo(clientId) // 从 API 获取游戏信息
	if err != nil {
		log.Fatalf("获取游戏信息失败: %v", err)
	}
	// 增加基础强度
	gameInfo.StrengthConfig.Strength -= PerfectQTEProgressBonus
	fmt.Printf("基础强度减小: %d\n", PerfectQTEProgressBonus)

	// 创建 SetStrengthConfigRequest 请求体
	config := http_api.SetStrengthConfigRequest{
		Strength: struct {
			Add *int `json:"add,omitempty"`
			Sub *int `json:"sub,omitempty"`
			Set *int `json:"set,omitempty"`
		}{
			Set: &gameInfo.StrengthConfig.Strength, // 传递增加的强度
		},
	}

	// 调用 http_api.SetStrength 更新基础强度
	err = http_api.SetStrength(clientId, config)
	if err != nil {
		log.Printf("更新基础强度失败: %v", err)
		return
	}
	fmt.Printf("成功更新基础强度: %d\n", gameInfo.StrengthConfig.Strength)
}
func missQTE() {

	gameInfo, err := http_api.GetGameInfo(clientId) // 从 API 获取游戏信息
	if err != nil {
		log.Fatalf("获取游戏信息失败: %v", err)
	}
	// 增加基础强度
	gameInfo.StrengthConfig.Strength += FailureProgressRetreat
	fmt.Printf("基础强度加: %d\n", FailureProgressRetreat)

	// 创建 SetStrengthConfigRequest 请求体
	config := http_api.SetStrengthConfigRequest{
		Strength: struct {
			Add *int `json:"add,omitempty"`
			Sub *int `json:"sub,omitempty"`
			Set *int `json:"set,omitempty"`
		}{
			Set: &gameInfo.StrengthConfig.Strength, // 传递增加的强度
		},
	}

	// 调用 http_api.SetStrength 更新基础强度
	err = http_api.SetStrength(clientId, config)
	if err != nil {
		log.Printf("更新基础强度失败: %v", err)
		return
	}
	fmt.Printf("成功更新基础强度: %d\n", gameInfo.StrengthConfig.Strength)
}
func main() {
	// Initialize window
	rl.InitWindow(800, 600, "QTE Game")
	defer rl.CloseWindow()

	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// QTE parameters
	x := float32(rl.GetScreenWidth()) / 2
	y := float32(rl.GetScreenHeight()) / 2
	radius := float32(150)
	pointerAngleStart := rand.Float32() * 360 // Random start angle between 0-360
	angleStart := rand.Float32() * 360
	direction := "clockwise" // Default rotation direction
	if rand.Float32() > 0.5 {
		direction = "counterclockwise" // Randomly select direction
	}

	// Variables for game state
	rotationActive := true
	var lastAngle float32
	var totalRotationAngle float32 // 累积旋转的角度
	var waitDuration time.Duration
	var waitEndTime time.Time
	judgeOnce := true

	// Main game loop
	for !rl.WindowShouldClose() {
		// Begin drawing
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Draw the circle
		rl.DrawCircleLines(int32(x), int32(y), radius, rl.DarkGray)

		// Draw the Great Skillcheck area (fixed arc)
		DrawArcPerfectQTE(rl.Vector2{X: x, Y: y}, radius-5, angleStart, angleStart+PerfectQTEAngle, 2, rl.Green)

		// Draw another arc (for demonstration) (fixed arc)
		DrawArcCommonQTE(rl.Vector2{X: x, Y: y}, radius-10, angleStart+PerfectQTEAngle, angleStart+PerfectQTEAngle+NormalQTEAngle, 2, rl.Blue)

		// Draw pointer rotation (this one rotates)
		if rotationActive {
			pointerAngleStart = DrawPointerRotation(rl.Vector2{X: x, Y: y}, radius, pointerAngleStart, direction)
			totalRotationAngle += rotationSpeedPerFrame // 累加旋转角度

			// Check if pointer has rotated 360 degrees
			if totalRotationAngle >= 360 {
				triggerFunction()      // 触发函数
				totalRotationAngle = 0 // 重置旋转角度
			}

		} else {
			DrawPointerRotation(rl.Vector2{X: x, Y: y}, radius, lastAngle, direction)
			// Check if the pointer landed in a specific QTE area when stopped
			if pointerAngleStart >= angleStart && pointerAngleStart <= angleStart+PerfectQTEAngle {
				if judgeOnce {
					rl.DrawText("Perfect QTE!", int32(x)-60, int32(y)-int32(radius)-40, 20, rl.Black)
					perfectQTE()
					judgeOnce = false
				}

			} else if pointerAngleStart >= angleStart+PerfectQTEAngle && pointerAngleStart <= angleStart+PerfectQTEAngle+NormalQTEAngle {
				if judgeOnce {
					rl.DrawText("Normal QTE!", int32(x)-60, int32(y)-int32(radius)-40, 20, rl.Black)
					judgeOnce = false
				}

			} else {
				if judgeOnce {
					rl.DrawText("Missed!", int32(x)-60, int32(y)-int32(radius)-40, 20, rl.Black)
					judgeOnce = false
					missQTE()
				}
			}

			// Random waiting logic
			if waitDuration == 0 {
				// Set random wait duration between 5 and 15 seconds
				waitDuration = time.Duration(rand.Int63n(int64(LongestWaitingTime)-int64(ShortestWaitingTime)) + int64(ShortestWaitingTime))
				waitEndTime = time.Now().Add(waitDuration)
			}

			// Check if the wait time has passed to restart the game
			if time.Now().After(waitEndTime) {
				// Restart the game with random parameters
				rotationActive = true
				pointerAngleStart = rand.Float32() * 360
				angleStart = rand.Float32() * 360
				direction = "clockwise" // Default rotation direction
				if rand.Float32() > 0.5 {
					direction = "counterclockwise" // Randomly select direction
				}
				waitDuration = 0 // Reset the wait duration
				judgeOnce = true
			}
		}

		// Draw the "Space/M1" button text
		rl.DrawText("Space/M1 to stop", int32(x)-40, int32(y)-int32(radius)-20, 20, rl.Black)

		// End drawing
		rl.EndDrawing()

		// Handle user input
		if rl.IsKeyDown(rl.KeySpace) || rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			rotationActive = false        // Stop the rotation when spacebar or mouse button is pressed
			lastAngle = pointerAngleStart // Store the angle where the pointer stopped
		}

		// Simulate the rotation frame rate
		time.Sleep(time.Second / frameRate)
	}
}
