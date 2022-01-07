package main

import (
	"github.com/verocity-gaming/unitehud/team"
)

// Configurations.
var (
	filenames = map[string]map[string][]filter{
		game: {
			gameVS: {
				filter{team.None, "img/game/vs.png", -0},
			},
			gameEnd: {
				filter{team.None, "img/game/end.png", -0},
			},
		},
		scored: {
			team.Purple.Name: {
				filter{team.Purple, "img/purple/score/score.png", -0},
				filter{team.Purple, "img/purple/score/score_alt.png", -0},
			},
			team.Orange.Name: {
				filter{team.Orange, "img/orange/score/score.png", -0},
			},
			team.Self.Name: {
				filter{team.Self, "img/self/score/score.png", -0},
				filter{team.Self, "img/self/score/score_alt.png", -0},
				filter{team.Self, "img/self/score/score_alt_alt.png", -0},
				filter{team.Self, "img/self/score/score_alt_alt_alt.png", -0},
				filter{team.Self, "img/self/score/score_alt_alt_alt_alt.png", -0},
				filter{team.Self, "img/self/score/score_alt_alt.png", -0},
				filter{team.Self, "img/self/score/score_big_alt.png", -0},
			},
		},
		points: {
			team.Purple.Name: {
				filter{team.Purple, "img/purple/points/point_0.png", 0},
				filter{team.Purple, "img/purple/points/point_0_alt.png", 0},
				filter{team.Purple, "img/purple/points/point_0_alt_alt.png", 0},
				filter{team.Purple, "img/purple/points/point_0_alt_alt_alt.png", 0},
				filter{team.Purple, "img/purple/points/point_0_alt_alt_alt_alt.png", 0},
				filter{team.Purple, "img/purple/points/point_0_alt_alt_alt_alt_alt.png", 0},
				filter{team.Purple, "img/purple/points/point_0_alt_alt_alt_alt_alt_alt.png", 0},
				filter{team.Purple, "img/purple/points/point_0_big.png", 0},
				filter{team.Purple, "img/purple/points/point_0_big_alt.png", 0},
				filter{team.Purple, "img/purple/points/point_0_big_alt_alt.png", 0},
				filter{team.Purple, "img/purple/points/point_0_big_alt_alt_alt.png", 0},
				filter{team.Purple, "img/purple/points/point_0_big_alt_alt_alt_alt.png", 0},

				filter{team.Purple, "img/purple/points/point_1.png", 1},
				filter{team.Purple, "img/purple/points/point_1_alt.png", 1},
				filter{team.Purple, "img/purple/points/point_1_alt_alt.png", 1},
				filter{team.Purple, "img/purple/points/point_1_big.png", 1},
				filter{team.Purple, "img/purple/points/point_1_big_alt.png", 1},
				filter{team.Purple, "img/purple/points/point_1_big_alt_alt.png", 1},

				filter{team.Purple, "img/purple/points/point_2.png", 2},
				filter{team.Purple, "img/purple/points/point_2_alt.png", 2},
				filter{team.Purple, "img/purple/points/point_2_alt_alt.png", 2},
				filter{team.Purple, "img/purple/points/point_2_alt_alt_alt.png", 2},
				filter{team.Purple, "img/purple/points/point_2_big_alt.png", 2},

				filter{team.Purple, "img/purple/points/point_3.png", 3},
				filter{team.Purple, "img/purple/points/point_3_alt.png", 3},

				filter{team.Purple, "img/purple/points/point_4.png", 4},
				filter{team.Purple, "img/purple/points/point_4_alt.png", 4},
				filter{team.Purple, "img/purple/points/point_4_alt_alt.png", 4},
				filter{team.Purple, "img/purple/points/point_4_big.png", 4},
				filter{team.Purple, "img/purple/points/point_4_big_alt.png", 4},
				filter{team.Purple, "img/purple/points/point_4_big_alt_alt.png", 4},
				filter{team.Purple, "img/purple/points/point_4_big_alt_alt_alt.png", 4},

				filter{team.Purple, "img/purple/points/point_5_alt.png", 5},
				filter{team.Purple, "img/purple/points/point_5_big.png", 5},

				filter{team.Purple, "img/purple/points/point_6.png", 6},
				filter{team.Purple, "img/purple/points/point_6_alt.png", 6},
				filter{team.Purple, "img/purple/points/point_6_big.png", 6},
				filter{team.Purple, "img/purple/points/point_6_big_alt.png", 6},

				filter{team.Purple, "img/purple/points/point_7.png", 7},
				filter{team.Purple, "img/purple/points/point_7_big.png", 7},

				filter{team.Purple, "img/purple/points/point_8.png", 8},
				filter{team.Purple, "img/purple/points/point_8_big.png", 8},
				filter{team.Purple, "img/purple/points/point_8_big_alt.png", 8},
				filter{team.Purple, "img/purple/points/point_8_big_alt_alt.png", 8},

				filter{team.Purple, "img/purple/points/point_9.png", 9},
				filter{team.Purple, "img/purple/points/point_9_alt.png", 9},
				filter{team.Purple, "img/purple/points/point_9_big.png", 9},
			},
			team.Orange.Name: {
				filter{team.Orange, "img/orange/points/point_0.png", 0},
				filter{team.Orange, "img/orange/points/point_0_alt.png", 0},
				filter{team.Orange, "img/orange/points/point_0_big.png", 0},
				filter{team.Orange, "img/orange/points/point_0_big_alt.png", 0},
				filter{team.Orange, "img/orange/points/point_0_big_alt_alt.png", 0},
				filter{team.Orange, "img/orange/points/point_0_big_alt_alt_alt.png", 0},
				filter{team.Orange, "img/orange/points/point_0_big_alt_alt_alt_alt.png", 0},

				filter{team.Orange, "img/orange/points/point_1.png", 1},
				filter{team.Orange, "img/orange/points/point_1_alt.png", 1},
				filter{team.Orange, "img/orange/points/point_1_big.png", 1},
				filter{team.Orange, "img/orange/points/point_1_big_alt.png", 1},

				filter{team.Orange, "img/orange/points/point_2.png", 2},
				filter{team.Orange, "img/orange/points/point_2_alt.png", 2},
				filter{team.Orange, "img/orange/points/point_2_big_alt.png", 2},

				filter{team.Orange, "img/orange/points/point_3.png", 3},
				filter{team.Orange, "img/orange/points/point_3_alt.png", 3},

				filter{team.Orange, "img/orange/points/point_4.png", 4},
				filter{team.Orange, "img/orange/points/point_4_alt.png", 4},
				filter{team.Orange, "img/orange/points/point_4_alt_alt.png", 4},
				filter{team.Orange, "img/orange/points/point_4_alt_alt_alt.png", 4},
				filter{team.Orange, "img/orange/points/point_4_big_alt.png", 4},

				filter{team.Orange, "img/orange/points/point_5.png", 5},
				filter{team.Orange, "img/orange/points/point_5_alt.png", 5},

				filter{team.Orange, "img/orange/points/point_6.png", 6},
				filter{team.Orange, "img/orange/points/point_6_alt.png", 6},
				filter{team.Orange, "img/orange/points/point_6_alt_alt.png", 6},
				filter{team.Orange, "img/orange/points/point_6_big_alt.png", 6},
				filter{team.Orange, "img/orange/points/point_6_big_alt_alt.png", 6},

				filter{team.Orange, "img/orange/points/point_7.png", 7},
				filter{team.Orange, "img/orange/points/point_7_big.png", 7},

				filter{team.Orange, "img/orange/points/point_8.png", 8},
				filter{team.Orange, "img/orange/points/point_8_alt.png", 8},
				filter{team.Orange, "img/orange/points/point_8_alt_alt.png", 8},
				filter{team.Orange, "img/orange/points/point_8_big_alt.png", 8},

				filter{team.Orange, "img/orange/points/point_9.png", 9},
				filter{team.Orange, "img/orange/points/point_9_alt.png", 9},
				filter{team.Orange, "img/orange/points/point_9_big.png", 9},
			},
			team.Self.Name: {
				filter{team.Self, "img/self/points/point_0.png", 0},
				filter{team.Self, "img/self/points/point_0_alt.png", 0},
				filter{team.Self, "img/self/points/point_0_alt_alt.png", 0},
				filter{team.Self, "img/self/points/point_0_alt_alt_alt.png", 0},
				filter{team.Self, "img/self/points/point_1.png", 1},
				filter{team.Self, "img/self/points/point_1_alt.png", 1},
				filter{team.Self, "img/self/points/point_5.png", 5},
				filter{team.Self, "img/self/points/point_5_alt.png", 5},
				filter{team.Self, "img/self/points/point_5_alt_alt.png", 5},
				filter{team.Self, "img/self/points/point_7.png", 7},
				filter{team.Self, "img/self/points/point_7_alt.png", 7},
				filter{team.Self, "img/self/points/point_7_alt_alt.png", 7},
				filter{team.Self, "img/self/points/point_8_alt.png", 8},
			},
		},
	}

	templates = map[string]map[string][]template{
		game: {
			team.None.Name: {},
		},
		scored: {
			team.Orange.Name: {},
			team.Purple.Name: {},
			team.Self.Name:   {},
		},
		points: {
			team.Orange.Name: {},
			team.Purple.Name: {},
			team.Self.Name:   {},
		},
	}
)
