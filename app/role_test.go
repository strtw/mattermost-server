// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"fmt"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/require"
)

func TestGetRolesByNames(t *testing.T) {
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

	channelScheme2, err := th.App.CreateScheme(&model.Scheme{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Scope:       model.SCHEME_SCOPE_CHANNEL,
	})
	require.Nil(t, err)
	defer th.App.DeleteScheme(channelScheme2.Id)
	channelSchemeRoles2, err := th.App.GetRolesByNames([]string{
		channelScheme2.DefaultChannelGuestRole,
		channelScheme2.DefaultChannelUserRole,
		channelScheme2.DefaultChannelAdminRole,
	})
	require.Nil(t, err)
	for _, role := range channelSchemeRoles2 {
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

	team2Channel2WithScheme.SchemeId = &channelScheme2.Id
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

	var channelGuest2, channelUser2, channelAdmin2 *model.Role

	for _, role := range channelSchemeRoles2 {
		switch role.Name {
		case channelScheme2.DefaultChannelGuestRole:
			channelGuest2 = role
		case channelScheme2.DefaultChannelUserRole:
			channelUser2 = role
		case channelScheme2.DefaultChannelAdminRole:
			channelAdmin2 = role
		}
	}

	fmt.Println(channelGuest2)
	fmt.Println(channelUser2)
	fmt.Println(channelAdmin2)

	cases := []struct {
		role         string
		permissionID string
		shouldHave   bool
	}{
		// The system roles keep the moderated and non-moderated permission
		{model.CHANNEL_ADMIN_ROLE_ID, nonModeratedPermission.Id, true},
		{model.CHANNEL_USER_ROLE_ID, nonModeratedPermission.Id, true},
		{model.CHANNEL_GUEST_ROLE_ID, nonModeratedPermission.Id, true},
		{model.CHANNEL_ADMIN_ROLE_ID, moderatedPermission.Id, true},
		{model.CHANNEL_USER_ROLE_ID, moderatedPermission.Id, true},
		{model.CHANNEL_GUEST_ROLE_ID, moderatedPermission.Id, true},

		// Reads from the channel scheme with a system scheme as the higher-scoped scheme
		{channelGuest.Name, moderatedPermission.Id, false},
		{channelUser.Name, moderatedPermission.Id, false},
		{channelAdmin.Name, moderatedPermission.Id, false},

		// Reads from the channel scheme with a team scheme as the higher-scoped scheme
		{channelGuest.Name, moderatedPermission.Id, false},
		{channelUser.Name, moderatedPermission.Id, false},
		{channelAdmin.Name, moderatedPermission.Id, false},

		// Reads from the system scheme for a non-moderated permission
		{channelGuest.Name, nonModeratedPermission.Id, true},
		{channelUser.Name, nonModeratedPermission.Id, true},
		{channelAdmin.Name, nonModeratedPermission.Id, true},

		// Reads from the system scheme for a non-moderated permission
		{channelGuest2.Name, permissionOnlyOnTeamScheme.Id, true},
		{channelUser2.Name, permissionOnlyOnTeamScheme.Id, true},
		{channelAdmin2.Name, permissionOnlyOnTeamScheme.Id, true},

		// Reads from the team scheme for a non-moderated permission
		{channelGuest.Name, permissionOnlyOnTeamScheme.Id, false},
		{channelUser.Name, permissionOnlyOnTeamScheme.Id, false},
		{channelAdmin.Name, permissionOnlyOnTeamScheme.Id, false},

		// Roles read from system scheme match those found on the channel
		{channelGuest.Name, guestOnlyPermission.Id, true},
		{channelUser.Name, guestOnlyPermission.Id, false},
		{channelAdmin.Name, guestOnlyPermission.Id, false},

		{channelGuest.Name, userOnlyPermission.Id, false},
		{channelUser.Name, userOnlyPermission.Id, true},
		{channelAdmin.Name, userOnlyPermission.Id, false},

		{channelGuest.Name, adminOnlyPermission.Id, false},
		{channelUser.Name, adminOnlyPermission.Id, false},
		{channelAdmin.Name, adminOnlyPermission.Id, true},

		// Roles read from team scheme match those found on the channel
		{channelGuest.Name, guestOnlyPermission.Id, true},
		{channelUser.Name, guestOnlyPermission.Id, false},
		{channelAdmin.Name, guestOnlyPermission.Id, false},

		{channelGuest.Name, userOnlyPermission.Id, false},
		{channelUser.Name, userOnlyPermission.Id, true},
		{channelAdmin.Name, userOnlyPermission.Id, false},

		{channelGuest.Name, adminOnlyPermission.Id, false},
		{channelUser.Name, adminOnlyPermission.Id, false},
		{channelAdmin.Name, adminOnlyPermission.Id, true},
	}

	roleNamePermissionsMap := map[string][]string{}
	allRoleNames := []string{}

	for _, testcase := range cases {
		// run through the tests once individually by role name
		modifiedRoles, err := th.App.GetRolesByNames([]string{testcase.role})
		require.Nil(t, err)
		require.Len(t, modifiedRoles, 1)
		role := modifiedRoles[0]
		require.Equal(t, testcase.shouldHave, includes(role.Permissions, testcase.permissionID), fmt.Sprintf("role.Name: %s, permission: %s", role.Name, testcase.permissionID))
		roleNamePermissionsMap[role.Name] = role.Permissions
		allRoleNames = append(allRoleNames, role.Name)
	}

	// verify that all the permission match when done in a batch
	modifiedRoles, err := th.App.GetRolesByNames(allRoleNames)
	require.Nil(t, err)
	for _, role := range modifiedRoles {
		require.ElementsMatch(t, roleNamePermissionsMap[role.Name], role.Permissions, fmt.Sprintf("role.Name: %s", role.Name))
	}
}
