package serverconfig

//func TestFlagsRESTORE_ENV_True(t *testing.T) {
//	t.Run("success", func(t *testing.T) {
//		err := os.Setenv("RESTORE", "true")
//		require.NoError(t, err)
//		flags, err := Flags()
//		require.NoError(t, err)
//		require.NotNil(t, flags.Restore)
//		require.Equal(t, flags.Restore, "true")
//	})
//}
//
//func TestFlagsRESTORE_ENV_FALSE(t *testing.T) {
//	t.Run("success", func(t *testing.T) {
//		err := os.Setenv("RESTORE", "false")
//		require.NoError(t, err)
//		flags, err := Flags()
//		require.NoError(t, err)
//		require.NotNil(t, flags.Restore)
//		require.Equal(t, flags.Restore, "false")
//	})
//}
//
//func TestFlagsRESTORE_ENV_EMPTY(t *testing.T) {
//	t.Run("success", func(t *testing.T) {
//		err := os.Setenv("RESTORE", "")
//		require.NoError(t, err)
//		flags, err := Flags()
//		require.NoError(t, err)
//		require.NotNil(t, flags.Restore)
//		require.Equal(t, flags.Restore, "true")
//	})
//}
