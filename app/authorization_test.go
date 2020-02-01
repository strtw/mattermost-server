// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/v5/model"
)

func TestCheckIfRolesGrantPermission(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	cases := []struct {
		roles        []string
		permissionId string
		shouldGrant  bool
	}{
		{[]string{model.SYSTEM_ADMIN_ROLE_ID}, model.PERMISSION_MANAGE_SYSTEM.Id, true},
		{[]string{model.SYSTEM_ADMIN_ROLE_ID}, "non-existent-permission", false},
		{[]string{model.CHANNEL_USER_ROLE_ID}, model.PERMISSION_READ_CHANNEL.Id, true},
		{[]string{model.CHANNEL_USER_ROLE_ID}, model.PERMISSION_MANAGE_SYSTEM.Id, false},
		{[]string{model.SYSTEM_ADMIN_ROLE_ID, model.CHANNEL_USER_ROLE_ID}, model.PERMISSION_MANAGE_SYSTEM.Id, true},
		{[]string{model.CHANNEL_USER_ROLE_ID, model.SYSTEM_ADMIN_ROLE_ID}, model.PERMISSION_MANAGE_SYSTEM.Id, true},
		{[]string{model.TEAM_USER_ROLE_ID, model.TEAM_ADMIN_ROLE_ID}, model.PERMISSION_MANAGE_SLASH_COMMANDS.Id, true},
		{[]string{model.TEAM_ADMIN_ROLE_ID, model.TEAM_USER_ROLE_ID}, model.PERMISSION_MANAGE_SLASH_COMMANDS.Id, true},
	}

	for _, testcase := range cases {
		assert.Equal(t, th.App.RolesGrantPermission(testcase.roles, testcase.permissionId), testcase.shouldGrant)
	}

}

func TestChannelRolesGrantPermission(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	moderatedPermission := model.PERMISSION_CREATE_POST
	nonModeratedPermission := model.PERMISSION_READ_CHANNEL
	permissionOnlyOnTeamScheme := model.PERMISSION_UPLOAD_FILE

	guestOnlyPermission := model.PERMISSION_EDIT_POST
	userOnlyPermission := model.PERMISSION_MANAGE_PUBLIC_CHANNEL_PROPERTIES
	adminOnlyPermission := model.PERMISSION_EDIT_OTHERS_POSTS

	th.App.SetLicense(model.NewTestLicense(""))
	th.App.SetPhase2PermissionsMigrationStatus(true)

	systemSchemeRoles, err := th.App.GetRolesByNames([]string{
		model.CHANNEL_GUEST_ROLE_ID,
		model.CHANNEL_USER_ROLE_ID,
		model.CHANNEL_ADMIN_ROLE_ID,
	})
	for _, role := range systemSchemeRoles {
		// defer resetting the system permissions back to what they were
		defer th.App.PatchRole(role, &model.RolePatch{
			Permissions: &role.Permissions,
		})

		permissions := []string{
			moderatedPermission.Id,
			nonModeratedPermission.Id,
		}

		if role.Name == model.CHANNEL_GUEST_ROLE_ID {
			permissions = append(permissions, guestOnlyPermission.Id)
		}
		if role.Name == model.CHANNEL_USER_ROLE_ID {
			permissions = append(permissions, userOnlyPermission.Id)
		}
		if role.Name == model.CHANNEL_ADMIN_ROLE_ID {
			permissions = append(permissions, adminOnlyPermission.Id)
		}

		_, err = th.App.PatchRole(role, &model.RolePatch{
			Permissions: &permissions,
		})
		require.Nil(t, err)
	}

	// make a team scheme with create_post and delete_post
	teamScheme, _ := th.CreateScheme()
	defer th.App.DeleteScheme(teamScheme.Id)
	teamSchemeRoles, err := th.App.GetRolesByNames([]string{
		teamScheme.DefaultChannelGuestRole,
		teamScheme.DefaultChannelUserRole,
		teamScheme.DefaultChannelAdminRole,
	})
	for _, role := range teamSchemeRoles {
		permissions := []string{
			moderatedPermission.Id,
			nonModeratedPermission.Id,
			permissionOnlyOnTeamScheme.Id,
		}

		if role.Name == teamScheme.DefaultChannelGuestRole {
			permissions = append(permissions, guestOnlyPermission.Id)
		}
		if role.Name == teamScheme.DefaultChannelUserRole {
			permissions = append(permissions, userOnlyPermission.Id)
		}
		if role.Name == teamScheme.DefaultChannelAdminRole {
			permissions = append(permissions, adminOnlyPermission.Id)
		}

		_, err = th.App.PatchRole(role, &model.RolePatch{
			Permissions: &permissions,
		})
		require.Nil(t, err)
	}

	// make a channel scheme without create_post and delete_post
	channelScheme, err := th.App.CreateScheme(&model.Scheme{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Scope:       model.SCHEME_SCOPE_CHANNEL,
	})
	require.Nil(t, err)
	defer th.App.DeleteScheme(channelScheme.Id)
	channelSchemeRoles, err := th.App.GetRolesByNames([]string{
		channelScheme.DefaultChannelGuestRole,
		channelScheme.DefaultChannelUserRole,
		channelScheme.DefaultChannelAdminRole,
	})
	require.Nil(t, err)
	for _, role := range channelSchemeRoles {
		_, err = th.App.PatchRole(role, &model.RolePatch{Permissions: &[]string{}})
		require.Nil(t, err)
	}

	// make two teams, one with the SchemeId set to the new team scheme
	team1 := th.CreateTeam()
	defer th.App.PermanentDeleteTeamId(team1.Id)

	team2WithScheme := th.CreateTeam()
	defer th.App.PermanentDeleteTeamId(team2WithScheme.Id)
	team2WithScheme.SchemeId = &teamScheme.Id
	team2WithScheme, err = th.App.UpdateTeamScheme(team2WithScheme)
	require.Nil(t, err)

	// make two channels per team, one per team with a channel scheme
	team1Channel1 := th.CreateChannel(team1)
	defer th.App.DeleteChannel(team1Channel1, "")

	team1Channel2WithScheme := th.CreateChannel(team1)
	defer th.App.DeleteChannel(team1Channel2WithScheme, "")

	team2Channel1 := th.CreateChannel(team2WithScheme)
	defer th.App.DeleteChannel(team2Channel1, "")

	team2Channel2WithScheme := th.CreateChannel(team2WithScheme)
	defer th.App.DeleteChannel(team2Channel2WithScheme, "")

	team1Channel2WithScheme.SchemeId = &channelScheme.Id
	team1Channel2WithScheme, err = th.App.UpdateChannel(team1Channel2WithScheme)
	require.Nil(t, err)

	team2Channel2WithScheme.SchemeId = &channelScheme.Id
	team2Channel2WithScheme, err = th.App.UpdateChannel(team2Channel2WithScheme)
	require.Nil(t, err)

	var channelGuest, channelUser, channelAdmin *model.Role

	for _, role := range channelSchemeRoles {
		switch role.Name {
		case channelScheme.DefaultChannelGuestRole:
			channelGuest = role
		case channelScheme.DefaultChannelUserRole:
			channelUser = role
		case channelScheme.DefaultChannelAdminRole:
			channelAdmin = role
		}
	}

	cases := []struct {
		roles        []string
		channelIDs   []string
		permissionID string
		shouldGrant  bool
	}{
		// The system roles grant the non-moderated permission to all of the channels
		{[]string{model.CHANNEL_ADMIN_ROLE_ID}, []string{team1Channel1.Id, team1Channel2WithScheme.Id, team2Channel1.Id, team2Channel2WithScheme.Id}, nonModeratedPermission.Id, true},
		{[]string{model.CHANNEL_USER_ROLE_ID}, []string{team1Channel1.Id, team1Channel2WithScheme.Id, team2Channel1.Id, team2Channel2WithScheme.Id}, nonModeratedPermission.Id, true},
		{[]string{model.CHANNEL_GUEST_ROLE_ID}, []string{team1Channel1.Id, team1Channel2WithScheme.Id, team2Channel1.Id, team2Channel2WithScheme.Id}, nonModeratedPermission.Id, true},

		// The system roles grant the non-moderated permission to all of the channels without schemes
		{[]string{model.CHANNEL_ADMIN_ROLE_ID}, []string{team1Channel1.Id, team2Channel1.Id}, moderatedPermission.Id, true},
		{[]string{model.CHANNEL_USER_ROLE_ID}, []string{team1Channel1.Id, team2Channel1.Id}, moderatedPermission.Id, true},
		{[]string{model.CHANNEL_GUEST_ROLE_ID}, []string{team1Channel1.Id, team2Channel1.Id}, moderatedPermission.Id, true},

		// The system roles do not grant the non-moderated permission to all of the channels with schemes
		{[]string{model.CHANNEL_ADMIN_ROLE_ID}, []string{team1Channel1.Id, team2Channel1.Id}, moderatedPermission.Id, true},
		{[]string{model.CHANNEL_USER_ROLE_ID}, []string{team1Channel1.Id, team2Channel1.Id}, moderatedPermission.Id, true},
		{[]string{model.CHANNEL_GUEST_ROLE_ID}, []string{team1Channel1.Id, team2Channel1.Id}, moderatedPermission.Id, true},

		// Reads from the channel scheme with a system scheme as the higher-scoped scheme
		{[]string{channelGuest.Name}, []string{team1Channel2WithScheme.Id}, moderatedPermission.Id, false},
		{[]string{channelUser.Name}, []string{team1Channel2WithScheme.Id}, moderatedPermission.Id, false},
		{[]string{channelAdmin.Name}, []string{team1Channel2WithScheme.Id}, moderatedPermission.Id, false},

		// Reads from the channel scheme with a team scheme as the higher-scoped scheme
		{[]string{channelGuest.Name}, []string{team2Channel2WithScheme.Id}, moderatedPermission.Id, false},
		{[]string{channelUser.Name}, []string{team2Channel2WithScheme.Id}, moderatedPermission.Id, false},
		{[]string{channelAdmin.Name}, []string{team2Channel2WithScheme.Id}, moderatedPermission.Id, false},

		// Reads from the system scheme for a non-moderated permission
		{[]string{channelGuest.Name}, []string{team1Channel2WithScheme.Id}, nonModeratedPermission.Id, true},
		{[]string{channelUser.Name}, []string{team1Channel2WithScheme.Id}, nonModeratedPermission.Id, true},
		{[]string{channelAdmin.Name}, []string{team1Channel2WithScheme.Id}, nonModeratedPermission.Id, true},

		// Reads from the team scheme for a non-moderated permission (setup check that it's false on the system scheme)
		{[]string{channelGuest.Name}, []string{team1Channel2WithScheme.Id}, permissionOnlyOnTeamScheme.Id, false},
		{[]string{channelUser.Name}, []string{team1Channel2WithScheme.Id}, permissionOnlyOnTeamScheme.Id, false},
		{[]string{channelAdmin.Name}, []string{team1Channel2WithScheme.Id}, permissionOnlyOnTeamScheme.Id, false},
		{[]string{channelGuest.Name}, []string{team2Channel2WithScheme.Id}, permissionOnlyOnTeamScheme.Id, true},
		{[]string{channelUser.Name}, []string{team2Channel2WithScheme.Id}, permissionOnlyOnTeamScheme.Id, true},
		{[]string{channelAdmin.Name}, []string{team2Channel2WithScheme.Id}, permissionOnlyOnTeamScheme.Id, true},

		// Roles read from system scheme match those found on the channel
		{[]string{channelGuest.Name}, []string{team1Channel2WithScheme.Id}, guestOnlyPermission.Id, true},
		{[]string{channelUser.Name}, []string{team1Channel2WithScheme.Id}, guestOnlyPermission.Id, false},
		{[]string{channelAdmin.Name}, []string{team1Channel2WithScheme.Id}, guestOnlyPermission.Id, false},

		{[]string{channelGuest.Name}, []string{team1Channel2WithScheme.Id}, userOnlyPermission.Id, false},
		{[]string{channelUser.Name}, []string{team1Channel2WithScheme.Id}, userOnlyPermission.Id, true},
		{[]string{channelAdmin.Name}, []string{team1Channel2WithScheme.Id}, userOnlyPermission.Id, false},

		{[]string{channelGuest.Name}, []string{team1Channel2WithScheme.Id}, adminOnlyPermission.Id, false},
		{[]string{channelUser.Name}, []string{team1Channel2WithScheme.Id}, adminOnlyPermission.Id, false},
		{[]string{channelAdmin.Name}, []string{team1Channel2WithScheme.Id}, adminOnlyPermission.Id, true},

		// Roles read from team scheme match those found on the channel
		{[]string{channelGuest.Name}, []string{team2Channel2WithScheme.Id}, guestOnlyPermission.Id, true},
		{[]string{channelUser.Name}, []string{team2Channel2WithScheme.Id}, guestOnlyPermission.Id, false},
		{[]string{channelAdmin.Name}, []string{team2Channel2WithScheme.Id}, guestOnlyPermission.Id, false},

		{[]string{channelGuest.Name}, []string{team2Channel2WithScheme.Id}, userOnlyPermission.Id, false},
		{[]string{channelUser.Name}, []string{team2Channel2WithScheme.Id}, userOnlyPermission.Id, true},
		{[]string{channelAdmin.Name}, []string{team2Channel2WithScheme.Id}, userOnlyPermission.Id, false},

		{[]string{channelGuest.Name}, []string{team2Channel2WithScheme.Id}, adminOnlyPermission.Id, false},
		{[]string{channelUser.Name}, []string{team2Channel2WithScheme.Id}, adminOnlyPermission.Id, false},
		{[]string{channelAdmin.Name}, []string{team2Channel2WithScheme.Id}, adminOnlyPermission.Id, true},
	}

	for _, testcase := range cases {
		for _, channelID := range testcase.channelIDs {
			assert.Equal(t, testcase.shouldGrant, th.App.ChannelRolesGrantPermission(testcase.roles, testcase.permissionID, channelID), fmt.Sprintf("roles: %+v, permission: %+v, channel: %s", testcase.roles, testcase.permissionID, channelID))
		}
	}

}
