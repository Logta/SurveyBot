package utils

func FindEmoji(num int) (string, error) {

	switch num {
	case 0:
		return "0️⃣", nil
	case 1:
		return "1️⃣", nil
	case 2:
		return "2️⃣", nil
	case 3:
		return "3️⃣", nil
	case 4:
		return "4️⃣", nil
	case 5:
		return "5️⃣", nil
	case 6:
		return "6️⃣", nil
	case 7:
		return "7️⃣", nil
	case 8:
		return "8️⃣", nil
	case 9:
		return "9️⃣", nil
	case 10:
		return "🔟", nil
	default:
		return "", errors.New("絵文字が見つかりません")
	}
}
