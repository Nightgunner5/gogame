package layout

var baseLayout = map[Coord]MultiTile{
	Coord{-6, -3}: {Wall1SE},
	Coord{-6, -2}: {Wall1},
	Coord{-6, -1}: {Wall1},
	Coord{-6, 0}:  {Wall1},
	Coord{-6, 1}:  {Wall1},
	Coord{-6, 2}:  {Wall1},
	Coord{-6, 3}:  {Wall1NE},
	Coord{-5, -3}: {Wall1},
	Coord{-5, -2}: {TileGray, Wall1NW},
	Coord{-5, -1}: {WireS, TileGray, Computer},
	Coord{-5, 0}:  {WireN, WireS, WireE, TileGray, Generator},
	Coord{-5, 1}:  {WireN, TileGray, Computer},
	Coord{-5, 2}:  {TileGray, Wall1SW},
	Coord{-5, 3}:  {Wall1},
	Coord{-4, -3}: {Wall1},
	Coord{-4, -2}: {TileGray},
	Coord{-4, -1}: {TileGray},
	Coord{-4, 0}:  {WireW, WireE, TileGray},
	Coord{-4, 1}:  {TileGray},
	Coord{-4, 2}:  {TileGray},
	Coord{-4, 3}:  {Wall1},
	Coord{-3, -3}: {Wall1},
	Coord{-3, -2}: {TileGray},
	Coord{-3, -1}: {TileGray},
	Coord{-3, 0}:  {WireW, WireE, TileGray},
	Coord{-3, 1}:  {TileGray},
	Coord{-3, 2}:  {TileGray},
	Coord{-3, 3}:  {TileBlack, DoorGeneralClosed},
	Coord{-2, -3}: {Wall1},
	Coord{-2, -2}: {TileGray},
	Coord{-2, -1}: {TileGray},
	Coord{-2, 0}:  {WireW, WireE, TileGray},
	Coord{-2, 1}:  {TileGray},
	Coord{-2, 2}:  {TileGray},
	Coord{-2, 3}:  {Wall1},
	Coord{-1, -4}: {Wall1SE},
	Coord{-1, -3}: {Wall1},
	Coord{-1, -2}: {TileWhite, Wall1},
	Coord{-1, -1}: {WireS, TileWhite, DoorGeneralClosed},
	Coord{-1, 0}:  {WireW, WireN, WireE, TileWhite, Wall1},
	Coord{-1, 1}:  {TileWhite, Window},
	Coord{-1, 2}:  {TileWhite, Wall1},
	Coord{-1, 3}:  {Wall1},
	Coord{-1, 4}:  {Wall1NE},
	Coord{-1, 7}:  {Window},
	Coord{-1, 8}:  {Window},
	Coord{-1, 9}:  {Wall1},
	Coord{-1, 10}: {Wall1},
	Coord{-1, 11}: {Wall1NE},
	Coord{0, -4}:  {Wall1},
	Coord{0, -3}:  {TileGray, Wall1NW},
	Coord{0, -2}:  {TileWhite, TileGrayNW},
	Coord{0, -1}:  {TileWhite, TriggerSelectRole},
	Coord{0, 0}:   {WireW, WireE, TileWhite},
	Coord{0, 1}:   {TileWhite, TriggerSelectRole},
	Coord{0, 2}:   {TileWhite, TileGraySW},
	Coord{0, 3}:   {TileGray, Wall1SW},
	Coord{0, 4}:   {Wall1},
	Coord{0, 7}:   {WireS, WireE, Window},
	Coord{0, 8}:   {WireN, WireS, TileRellow, TileGraySE},
	Coord{0, 9}:   {WireN, WireS, TileRellow, TileGraySE, TileGrayNE},
	Coord{0, 10}:  {WireN, WireS, TileRellow, TileGraySE, TileGrayNE, Wall1SW},
	Coord{0, 11}:  {WireN, WireS, TileRellow, Wall1},
	Coord{0, 12}:  {WireN, WireE, WireS, TileRellow, Wall1},
	Coord{0, 13}:  {WireN, WireS, TileRellow, Wall1},
	Coord{0, 14}:  {WireN, WireS, Window},
	Coord{0, 15}:  {WireE, WireN, Window},
	Coord{1, -4}:  {Wall1},
	Coord{1, -3}:  {WireN, WireS, TileWhite, Light1NOn},
	Coord{1, -2}:  {WireN, WireS, TileWhite},
	Coord{1, -1}:  {WireN, WireS, TileWhite},
	Coord{1, 0}:   {WireW, WireN, WireS, WireE, TileWhite, TriggerSelectRole},
	Coord{1, 1}:   {WireN, WireS, TileWhite},
	Coord{1, 2}:   {WireN, WireS, TileWhite},
	Coord{1, 3}:   {WireS, WireN, TileWhite, Light1SOn},
	Coord{1, 4}:   {Wall1},
	Coord{1, 7}:   {WireW, WireE, Window},
	Coord{1, 8}:   {TileRellow, TileGraySE, TileGraySW},
	Coord{1, 9}:   {TileGray},
	Coord{1, 10}:  {TileGray},
	Coord{1, 11}:  {TileGray},
	Coord{1, 12}:  {TileGray, Light1WOn},
	Coord{1, 13}:  {TileGray},
	Coord{1, 14}:  {TileGray},
	Coord{1, 15}:  {WireW, WireE, Window},
	Coord{2, -4}:  {Wall1},
	Coord{2, -3}:  {TileWhite, TileGrayNE, TileGrayNW},
	Coord{2, -2}:  {TileWhite},
	Coord{2, -1}:  {TileWhite},
	Coord{2, 0}:   {WireW, WireE, TileWhite},
	Coord{2, 1}:   {TileWhite},
	Coord{2, 2}:   {TileWhite},
	Coord{2, 3}:   {TileWhite, TileGraySW, TileGraySE},
	Coord{2, 4}:   {Wall1},
	Coord{2, 5}:   {Window},
	Coord{2, 6}:   {Window},
	Coord{2, 7}:   {WireW, WireE, Wall1},
	Coord{2, 8}:   {WireE, TileGray, Light1NOn},
	Coord{2, 9}:   {TileGray},
	Coord{2, 10}:  {TileGray},
	Coord{2, 11}:  {TileGray},
	Coord{2, 12}:  {TileGray},
	Coord{2, 13}:  {TileGray},
	Coord{2, 14}:  {TileGray},
	Coord{2, 15}:  {WireW, WireE, Window},
	Coord{3, -4}:  {Wall1},
	Coord{3, -3}:  {TileWhite, TileGrayNE, TileGrayNW},
	Coord{3, -2}:  {TileWhite},
	Coord{3, -1}:  {TileWhite},
	Coord{3, 0}:   {WireW, WireE, TileWhite},
	Coord{3, 1}:   {TileWhite},
	Coord{3, 2}:   {TileWhite},
	Coord{3, 3}:   {TileWhite, TileGraySW, TileGraySE},
	Coord{3, 4}:   {WireE, TileBlack, DoorGeneralClosed},
	Coord{3, 5}:   {TileBlack},
	Coord{3, 6}:   {TileBlack},
	Coord{3, 7}:   {WireW, TileBlack, DoorGeneralClosed},
	Coord{3, 8}:   {WireE, WireW, TileGray},
	Coord{3, 9}:   {TileGray},
	Coord{3, 10}:  {TileGray},
	Coord{3, 11}:  {TileGray},
	Coord{3, 12}:  {TileGray},
	Coord{3, 13}:  {TileGray},
	Coord{3, 14}:  {TileGray},
	Coord{3, 15}:  {WireW, WireE, Window},
	Coord{4, -4}:  {Wall1},
	Coord{4, -3}:  {WireN, WireS, TileWhite, Light1NOn},
	Coord{4, -2}:  {WireN, WireS, TileWhite},
	Coord{4, -1}:  {WireN, WireS, TileWhite},
	Coord{4, 0}:   {WireW, WireN, WireS, TileWhite},
	Coord{4, 1}:   {WireN, WireS, TileWhite},
	Coord{4, 2}:   {WireN, WireS, TileWhite},
	Coord{4, 3}:   {WireN, WireS, TileWhite, Light1SOn},
	Coord{4, 4}:   {WireN, WireW, Wall1},
	Coord{4, 5}:   {Window},
	Coord{4, 6}:   {Window},
	Coord{4, 7}:   {Wall1},
	Coord{4, 8}:   {WireE, WireW, TileGray, Light1NOn},
	Coord{4, 9}:   {TileGray},
	Coord{4, 10}:  {TileGray},
	Coord{4, 11}:  {TileGray},
	Coord{4, 12}:  {TileGray},
	Coord{4, 13}:  {TileGray},
	Coord{4, 14}:  {TileGray},
	Coord{4, 15}:  {WireE, WireW, Window},
	Coord{5, -4}:  {Wall1},
	Coord{5, -3}:  {TileGray, Wall1NE},
	Coord{5, -2}:  {TileWhite, TileGrayNE},
	Coord{5, -1}:  {TileWhite},
	Coord{5, 0}:   {TileWhite},
	Coord{5, 1}:   {TileWhite},
	Coord{5, 2}:   {TileWhite, TileGraySE},
	Coord{5, 3}:   {TileGray, Wall1SE},
	Coord{5, 4}:   {Wall1},
	Coord{5, 7}:   {Window},
	Coord{5, 8}:   {WireE, WireW, TileRellow, TileGraySE, TileGraySW},
	Coord{5, 9}:   {TileGray},
	Coord{5, 10}:  {TileGray},
	Coord{5, 11}:  {TileGray},
	Coord{5, 12}:  {TileGray},
	Coord{5, 13}:  {TileGray},
	Coord{5, 14}:  {TileGray},
	Coord{5, 15}:  {WireS, WireW, Wall1},
	Coord{5, 16}:  {WireE, WireN, Wall1},
	Coord{5, 17}:  {Wall1NE},
	Coord{6, -4}:  {Wall1SW},
	Coord{6, -3}:  {Wall1},
	Coord{6, -2}:  {TileGray, Wall1NE},
	Coord{6, -1}:  {TileWhite, TileGrayNE},
	Coord{6, 0}:   {TileWhite},
	Coord{6, 1}:   {TileWhite, TileGraySE},
	Coord{6, 2}:   {TileGray, Wall1SE},
	Coord{6, 3}:   {Wall1},
	Coord{6, 4}:   {Wall1NW},
	Coord{6, 7}:   {Window},
	Coord{6, 8}:   {WireE, WireW, TileRellow, TileGraySE, TileGraySW},
	Coord{6, 9}:   {TileGray},
	Coord{6, 10}:  {TileGray},
	Coord{6, 11}:  {TileGray, TileRellowSE},
	Coord{6, 12}:  {TileRellow, TileGrayNW, TileGraySW, Light1EOn},
	Coord{6, 13}:  {TileRellow, TileGrayNW, TileGraySW},
	Coord{6, 14}:  {TileGray, TileRellowSE, TileRellowNE},
	Coord{6, 15}:  {TileGray, TileRellowSE, TileRellowNE, Window},
	Coord{6, 16}:  {WireE, WireW, TileRellow, Generator},
	Coord{6, 17}:  {Wall1},
	Coord{6, 18}:  {Wall1NE},
	Coord{7, -3}:  {Wall1SW},
	Coord{7, -2}:  {Wall1},
	Coord{7, -1}:  {Window},
	Coord{7, 0}:   {Window},
	Coord{7, 1}:   {Window},
	Coord{7, 2}:   {Wall1},
	Coord{7, 3}:   {Wall1NW},
	Coord{7, 7}:   {Wall1},
	Coord{7, 8}:   {WireE, WireW, TileRellow, TileGraySE, TileGraySW, Wall1NE},
	Coord{7, 9}:   {TileGray},
	Coord{7, 10}:  {TileGray},
	Coord{7, 11}:  {TileRellow, TileGrayNW, TileGrayNE, Wall1SE},
	Coord{7, 12}:  {WireE, WireW, TileRellow, Wall1},
	Coord{7, 13}:  {WireS, TileRellow, DoorGeneralClosed},
	Coord{7, 14}:  {WireS, WireN, TileRellow, Wall1},
	Coord{7, 15}:  {WireS, WireN, Wall1},
	Coord{7, 16}:  {WireN, WireW, Wall1},
	Coord{7, 17}:  {Wall1},
	Coord{7, 18}:  {Wall1},
	Coord{7, 19}:  {Wall1},
	Coord{7, 20}:  {Wall1},
	Coord{7, 21}:  {Wall1},
	Coord{7, 22}:  {Wall1},
	Coord{7, 23}:  {Wall1},
	Coord{7, 24}:  {Wall1},
	Coord{7, 25}:  {Wall1},
	Coord{8, 7}:   {Wall1},
	Coord{8, 8}:   {WireE, WireW, Wall1},
	Coord{8, 9}:   {TileBlack, Window},
	Coord{8, 10}:  {TileBlack, Window},
	Coord{8, 11}:  {TileBlack, Wall1},
	Coord{8, 12}:  {WireE, WireW, TileGray, TileWhiteSE, Wall1NW},
	Coord{8, 13}:  {TileWhite},
	Coord{8, 14}:  {TileWhite, TileGraySW, TileGrayNW},
	Coord{8, 15}:  {TileWhite, TileGraySW, TileGrayNW},
	Coord{8, 16}:  {TileWhite, TileGraySW, TileGrayNW},
	Coord{8, 17}:  {TileWhite, TileGraySW, TileGrayNW},
	Coord{8, 18}:  {TileWhite, TileGraySW, TileGrayNW, Light1WOn},
	Coord{8, 19}:  {TileWhite, TileGraySW, TileGrayNW},
	Coord{8, 20}:  {TileGray, TileWhiteNE, Wall1SW},
	Coord{8, 21}:  {Wall1},
	Coord{8, 22}:  {WireE, WireW, WireS, TileGray, 46},
	Coord{8, 23}:  {WireN, WireS, TileGray, 46},
	Coord{8, 24}:  {WireN, TileGray, 46},
	Coord{8, 25}:  {Wall1},
	Coord{9, 7}:   {Wall1},
	Coord{9, 8}:   {WireE, WireW, TileRed, Safe},
	Coord{9, 9}:   {TileBlack},
	Coord{9, 10}:  {TileBlack},
	Coord{9, 11}:  {TileBlack, Window},
	Coord{9, 12}:  {WireE, WireW, TileWhite, TileGrayNW, TileGrayNE},
	Coord{9, 13}:  {TileWhite},
	Coord{9, 14}:  {TileWhite},
	Coord{9, 15}:  {TileWhite},
	Coord{9, 16}:  {TileWhite},
	Coord{9, 17}:  {TileWhite},
	Coord{9, 18}:  {TileWhite},
	Coord{9, 19}:  {TileWhite},
	Coord{9, 20}:  {TileWhite, TileGraySE, TileGraySW},
	Coord{9, 21}:  {Wall1},
	Coord{9, 22}:  {WireE, WireW, TileWhite},
	Coord{9, 23}:  {TileWhite},
	Coord{9, 24}:  {TileWhite},
	Coord{9, 25}:  {Wall1},
	Coord{10, 7}:  {Wall1},
	Coord{10, 8}:  {WireW, WireE, TileRed, Computer},
	Coord{10, 9}:  {TileBlack},
	Coord{10, 10}: {TileBlack},
	Coord{10, 11}: {TileBlack, Window},
	Coord{10, 12}: {WireW, WireS, TileWhite, TileGrayNW, TileGrayNE},
	Coord{10, 13}: {WireS, WireN, TileWhite},
	Coord{10, 14}: {WireN, WireE, TileWhite},
	Coord{10, 15}: {TileWhite},
	Coord{10, 16}: {TileWhite},
	Coord{10, 17}: {TileWhite},
	Coord{10, 18}: {TileWhite},
	Coord{10, 19}: {TileWhite},
	Coord{10, 20}: {TileWhite, TileGraySE, TileGraySW},
	Coord{10, 21}: {Wall1},
	Coord{10, 22}: {WireW, WireE, TileWhite},
	Coord{10, 23}: {TileWhite},
	Coord{10, 24}: {TileWhite},
	Coord{10, 25}: {Wall1},
	Coord{11, 7}:  {Wall1},
	Coord{11, 8}:  {WireS, WireW, TileRed, Computer},
	Coord{11, 9}:  {WireE, WireS, WireN, TileBlack, Light1EOn},
	Coord{11, 10}: {WireE, WireN, TileBlack},
	Coord{11, 11}: {TileBlack, Wall1},
	Coord{11, 12}: {TileGray, TileWhiteSW, Wall1NE},
	Coord{11, 13}: {TileWhite, TileGrayNE, TileGraySE},
	Coord{11, 14}: {WireE, WireW, TileWhite, TileGrayNE, TileGraySE, Light1EOn},
	Coord{11, 15}: {TileWhite, TileGrayNE, TileGraySE},
	Coord{11, 16}: {TileWhite, TileGrayNE},
	Coord{11, 17}: {TileWhite},
	Coord{11, 18}: {TileWhite},
	Coord{11, 19}: {TileWhite},
	Coord{11, 20}: {TileWhite, TileGraySE, TileGraySW},
	Coord{11, 21}: {Wall1},
	Coord{11, 22}: {WireW, WireE, TileWhite},
	Coord{11, 23}: {TileWhite},
	Coord{11, 24}: {TileWhite},
	Coord{11, 25}: {Wall1},
	Coord{12, 7}:  {Wall1SW},
	Coord{12, 8}:  {Wall1},
	Coord{12, 9}:  {TileBlack, Wall1},
	Coord{12, 10}: {WireE, WireW, TileBlack, DoorSecurityClosed},
	Coord{12, 11}: {TileBlack, Wall1},
	Coord{12, 12}: {TileBlack, Wall1},
	Coord{12, 13}: {TileBlack, Wall1},
	Coord{12, 14}: {WireE, WireW, TileBlack, Wall1},
	Coord{12, 15}: {TileBlack, Wall1},
	Coord{12, 16}: {TileGray, Wall1NE},
	Coord{12, 17}: {TileWhite, TileGrayNE},
	Coord{12, 18}: {TileWhite},
	Coord{12, 19}: {TileWhite},
	Coord{12, 20}: {TileWhite, TileGraySE, TileGraySW},
	Coord{12, 21}: {Wall1},
	Coord{12, 22}: {WireW, WireE, TileWhite, TileWhite},
	Coord{12, 23}: {TileWhite},
	Coord{12, 24}: {TileWhite},
	Coord{12, 25}: {Wall1},
	Coord{13, 8}:  {Wall1},
	Coord{13, 9}:  {WireE, WireS, TileBlack, Wall1NW},
	Coord{13, 10}: {WireW, WireN, TileBlack},
	Coord{13, 11}: {TileBlack},
	Coord{13, 12}: {WireE, TileRed, Light1WOn},
	Coord{13, 13}: {TileRed},
	Coord{13, 14}: {WireE, WireW, TileRed},
	Coord{13, 15}: {TileRed, Wall1SW},
	Coord{13, 16}: {TileBlack, Wall1},
	Coord{13, 17}: {WireN, WireE, TileBlack, Light1NOn},
	Coord{13, 18}: {TileWhite},
	Coord{13, 19}: {TileWhite},
	Coord{13, 20}: {WireE, WireS, TileRellow, Light1SOn},
	Coord{13, 21}: {WireS, WireN, Wall1},
	Coord{13, 22}: {WireW, WireE, WireN, TileWhite},
	Coord{13, 23}: {TileWhite},
	Coord{13, 24}: {TileWhite},
	Coord{13, 25}: {Wall1},
	Coord{14, 8}:  {Wall1},
	Coord{14, 9}:  {WireW, WireS, TileBlack, Generator},
	Coord{14, 10}: {WireN, WireS, TileBlack},
	Coord{14, 11}: {WireN, WireS, WireE, TileBlack},
	Coord{14, 12}: {WireN, WireS, WireW, TileBlack},
	Coord{14, 13}: {WireN, WireE, WireS, TileBlack},
	Coord{14, 14}: {WireN, WireS, WireW, TileBlack},
	Coord{14, 15}: {WireN, WireS, TileBlack},
	Coord{14, 16}: {WireN, WireS, TileBlack, DoorSecurityClosed},
	Coord{14, 17}: {WireS, WireN, WireE, WireW, TileBlack},
	Coord{14, 18}: {WireS, WireN, TileWhite},
	Coord{14, 19}: {WireS, WireN, TileWhite},
	Coord{14, 20}: {WireN, WireE, WireW, TileRellow},
	Coord{14, 21}: {WireE, TileRellow, DoorEngineerClosed},
	Coord{14, 22}: {WireW, WireN, TileWhite},
	Coord{14, 23}: {TileWhite},
	Coord{14, 24}: {TileWhite},
	Coord{14, 25}: {Wall1},
	Coord{15, 8}:  {Wall1},
	Coord{15, 9}:  {TileBlack, Wall1NE},
	Coord{15, 10}: {TileBlack},
	Coord{15, 11}: {WireE, WireW, TileBlack},
	Coord{15, 12}: {TileRed},
	Coord{15, 13}: {WireW, TileRed, Computer},
	Coord{15, 14}: {TileRed},
	Coord{15, 15}: {TileRed, Wall1SE},
	Coord{15, 16}: {TileBlack, Wall1},
	Coord{15, 17}: {WireN, WireW, TileBlack, Light1NOn},
	Coord{15, 18}: {TileWhite},
	Coord{15, 19}: {TileWhite},
	Coord{15, 20}: {TileRellow, Light1SOn},
	Coord{15, 21}: {WireE, WireW, Wall1},
	Coord{15, 22}: {TileWhite},
	Coord{15, 23}: {TileWhite},
	Coord{15, 24}: {TileWhite},
	Coord{15, 25}: {Wall1},
	Coord{16, 8}:  {Wall1SW},
	Coord{16, 9}:  {Wall1},
	Coord{16, 10}: {Wall1},
	Coord{16, 11}: {WireE, WireW, TileBlack},
	Coord{16, 12}: {TileBlack, Wall1},
	Coord{16, 13}: {TileBlack, Wall1},
	Coord{16, 14}: {TileBlack, Wall1},
	Coord{16, 15}: {TileBlack, Wall1},
	Coord{16, 16}: {TileBlack, Wall1},
	Coord{16, 17}: {TileBlack, Wall1NE},
	Coord{16, 18}: {TileWhite},
	Coord{16, 19}: {TileWhite},
	Coord{16, 20}: {TileRellow, Wall1SE},
	Coord{16, 21}: {WireS, WireW, Wall1},
	Coord{16, 22}: {WireN, TileWhite, Generator},
	Coord{16, 23}: {TileWhite},
	Coord{16, 24}: {TileWhite},
	Coord{16, 25}: {Wall1},
	Coord{17, 9}:  {Wall1},
	Coord{17, 10}: {Wall1},
	Coord{17, 11}: {WireE, WireW, TileBlack},
	Coord{17, 12}: {TileBlack, Wall1},
	Coord{17, 13}: {TileGray},
	Coord{17, 14}: {WireE, TileGray, Light1WOn},
	Coord{17, 15}: {TileGray},
	Coord{17, 16}: {TileGray},
	Coord{17, 17}: {TileWhite, Wall1},
	Coord{17, 18}: {TileWhite},
	Coord{17, 19}: {TileWhite, TileCyanSE},
	Coord{17, 20}: {TileWhite, Wall1},
	Coord{17, 21}: {Wall1},
	Coord{17, 22}: {Wall1},
	Coord{17, 23}: {Wall1},
	Coord{17, 24}: {Wall1},
	Coord{17, 25}: {Wall1},
	Coord{18, 9}:  {Wall1},
	Coord{18, 10}: {WireS, TileBlack, Computer},
	Coord{18, 11}: {WireE, WireN, WireW, TileBlack},
	Coord{18, 12}: {TileBlack, Wall1},
	Coord{18, 13}: {TileGray},
	Coord{18, 14}: {WireW, WireE, TileGray},
	Coord{18, 15}: {TileGray},
	Coord{18, 16}: {TileGray},
	Coord{18, 17}: {TileWhite, Window},
	Coord{18, 18}: {TileWhite},
	Coord{18, 19}: {TileWhite, TileCyanSE, TileCyanNW},
	Coord{18, 20}: {TileWhite, Window},
	Coord{18, 21}: {TileGray},
	Coord{18, 22}: {TileGray},
	Coord{18, 23}: {Wall1},
	Coord{18, 24}: {Wall1},
	Coord{18, 25}: {Wall1NW},
	Coord{19, 9}:  {Wall1SW},
	Coord{19, 10}: {Wall1},
	Coord{19, 11}: {WireW, WireS, TileBlack},
	Coord{19, 12}: {WireN, WireS, TileBlack, DoorSecurityClosed},
	Coord{19, 13}: {WireN, WireS, TileGray},
	Coord{19, 14}: {WireN, WireW, WireE, TileGray},
	Coord{19, 15}: {TileGray},
	Coord{19, 16}: {TileGray},
	Coord{19, 17}: {TileWhite, Window},
	Coord{19, 18}: {TileWhite},
	Coord{19, 19}: {TileWhite, TileCyanSE, TileCyanNW},
	Coord{19, 20}: {TileWhite, Window},
	Coord{19, 21}: {TileGray},
	Coord{19, 22}: {TileGray},
	Coord{19, 23}: {Wall1},
	Coord{19, 24}: {Wall1NW},
	Coord{20, 10}: {Wall1SW},
	Coord{20, 11}: {TileBlack, Wall1},
	Coord{20, 12}: {TileBlack, Wall1},
	Coord{20, 13}: {TileGray},
	Coord{20, 14}: {WireW, WireE, TileGray},
	Coord{20, 15}: {TileGray},
	Coord{20, 16}: {TileGray},
	Coord{20, 17}: {TileWhite, Wall1},
	Coord{20, 18}: {TileWhite},
	Coord{20, 19}: {WireS, WireE, TileWhite, TileCyanSE, TileCyanNW, Light1SOn},
	Coord{20, 20}: {TileWhite, Wall1},
	Coord{20, 21}: {TileGray},
	Coord{20, 22}: {TileGray},
	Coord{20, 23}: {Wall1},
	Coord{21, 11}: {Wall1SW},
	Coord{21, 12}: {TileBlack, Wall1},
	Coord{21, 13}: {TileGray},
	Coord{21, 14}: {WireW, WireS, WireE, TileGray, Light1EOn},
	Coord{21, 15}: {WireN, WireE, WireS, TileGray},
	Coord{21, 16}: {WireN, WireS, TileGray},
	Coord{21, 17}: {WireN, WireE, TileWhite, DoorMedicalClosed},
	Coord{21, 18}: {TileWhite},
	Coord{21, 19}: {WireE, WireW, WireS, TileWhite, TileCyanNW, TileCyanSE},
	Coord{21, 20}: {WireN, TileWhite, DoorMedicalClosed},
	Coord{21, 21}: {TileGray},
	Coord{21, 22}: {TileGray, Light1EOn},
	Coord{21, 23}: {Wall1},
	Coord{22, 12}: {Wall1SW},
	Coord{22, 13}: {TileWhite, Wall1},
	Coord{22, 14}: {WireE, WireW, TileWhite, Wall1},
	Coord{22, 15}: {WireW, TileWhite, DoorMedicalClosed},
	Coord{22, 16}: {TileWhite, Wall1},
	Coord{22, 17}: {WireW, WireE, WireS, TileWhite, Wall1},
	Coord{22, 18}: {WireN, WireS, TileWhite, DoorMedicalOpen},
	Coord{22, 19}: {WireS, WireN, WireW, TileWhite, Window},
	Coord{22, 20}: {WireN, WireE, WireS, TileWhite, Wall1},
	Coord{22, 21}: {WireN, WireS, TileWhite, DoorMedicalClosed},
	Coord{22, 22}: {WireN, WireE, WireW, TileWhite, Wall1},
	Coord{22, 23}: {Wall1},
	Coord{23, 13}: {Wall1},
	Coord{23, 14}: {WireW, WireE, TileWhite},
	Coord{23, 15}: {TileWhite},
	Coord{23, 16}: {TileWhite},
	Coord{23, 17}: {TileWhite, Light1WOn},
	Coord{23, 18}: {TileWhite},
	Coord{23, 19}: {TileWhite},
	Coord{23, 20}: {TileWhite, Light1WOn},
	Coord{23, 21}: {TileWhite},
	Coord{23, 22}: {WireW, WireE, TileWhite},
	Coord{23, 23}: {Wall1},
	Coord{24, 13}: {Wall1},
	Coord{24, 14}: {WireN, WireW, TileWhite, Light1NOn},
	Coord{24, 15}: {TileWhite},
	Coord{24, 16}: {TileWhite},
	Coord{24, 17}: {TileWhite},
	Coord{24, 18}: {TileWhite},
	Coord{24, 19}: {TileWhite},
	Coord{24, 20}: {TileWhite},
	Coord{24, 21}: {TileWhite},
	Coord{24, 22}: {WireS, WireW, TileWhite, Light1SOn},
	Coord{24, 23}: {Wall1},
	Coord{25, 13}: {Wall1},
	Coord{25, 14}: {TileWhite},
	Coord{25, 15}: {TileWhite},
	Coord{25, 16}: {TileWhite},
	Coord{25, 17}: {TileWhite},
	Coord{25, 18}: {TileWhite},
	Coord{25, 19}: {TileWhite},
	Coord{25, 20}: {TileWhite},
	Coord{25, 21}: {TileWhite},
	Coord{25, 22}: {TileWhite},
	Coord{25, 23}: {Wall1},
	Coord{26, 13}: {Wall1SW},
	Coord{26, 14}: {Wall1},
	Coord{26, 15}: {Wall1},
	Coord{26, 16}: {Wall1},
	Coord{26, 17}: {Wall1},
	Coord{26, 18}: {Wall1},
	Coord{26, 19}: {Wall1},
	Coord{26, 20}: {Wall1},
	Coord{26, 21}: {Wall1},
	Coord{26, 22}: {Wall1},
	Coord{26, 23}: {Wall1NW},
}
