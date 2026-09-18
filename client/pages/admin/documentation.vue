<script lang="ts" setup>
definePageMeta({
  layout: 'loggedin'
})
</script>

<template>
  <div class="bg-gray-50 min-h-screen py-8 px-4 sm:px-6 lg:px-8">
    <div class="max-w-4xl mx-auto space-y-10">

      <div class="flex flex-wrap items-center justify-between gap-4">
        <h1 class="text-3xl font-bold text-gray-900">Documentation administrateur</h1>
      </div>

      <p class="text-gray-600">
        Cette page recense toutes les actions gérables depuis l'interface d'administration.
      </p>

      <section class="bg-white rounded-2xl shadow p-8 space-y-4">
        <h2 class="text-xl font-semibold text-gray-900">🔑 Rôles & permissions</h2>
        <p class="text-gray-600">
          Chaque compte peut cumuler trois rôles indépendants :
        </p>
        <ul class="list-disc list-inside text-gray-600 space-y-1">
          <li><strong>Tuteur</strong> : peut s'inscrire pour encadrer une matière et déclarer ses séances/heures.</li>
          <li><strong>Tutee</strong> : peut s'inscrire à une matière pour être accompagné.</li>
          <li><strong>Administrateur</strong> : accède à cet espace, gère campagnes, matières, utilisateurs et affectations.</li>
        </ul>
        <p class="text-gray-600">
          Pour modifier les rôles d'un utilisateur, rendez-vous dans
          <NuxtLink to="/admin/users" class="text-blue-600 hover:underline">Admin &gt; Utilisateurs</NuxtLink>,
          puis cliquez sur « Modifier » sur la ligne correspondante. Par sécurité, un administrateur ne peut pas
          retirer ses propres droits admin (pour éviter de se retrouver bloqué), un autre admin doit s'en charger.
        </p>
      </section>

      <section class="bg-white rounded-2xl shadow p-8 space-y-4">
        <h2 class="text-xl font-semibold text-gray-900">📅 Campagnes</h2>
        <p class="text-gray-600">
          Depuis <NuxtLink to="/admin/campaigns" class="text-blue-600 hover:underline">Admin &gt; Campagnes</NuxtLink>,
          vous pouvez créer une nouvelle campagne de tutorat (« + Nouvelle campagne ») en renseignant l'année
          scolaire, le semestre, les dates de la campagne et la fenêtre d'inscription. Cliquez sur « Voir » pour
          modifier une campagne existante ou consulter son détail.
        </p>
        <p class="text-gray-600">
          Une campagne terminée peut être <strong>archivée</strong> (bouton « Archiver ») plutôt que supprimée :
          elle disparaît des tableaux de bord des tuteurs et tutorés (ils ne la voient plus, ne peuvent plus s'y
          inscrire), mais reste entièrement consultable et modifiable depuis l'admin,  ses données (affectations,
          heures, séances) ne sont jamais perdues. L'action est réversible via « Désarchiver ». Il n'existe
          volontairement aucune suppression définitive de campagne.
        </p>
      </section>

      <section class="bg-white rounded-2xl shadow p-8 space-y-4">
        <h2 class="text-xl font-semibold text-gray-900">📚 Matières</h2>
        <p class="text-gray-600">
          Depuis <NuxtLink to="/admin/subjects" class="text-blue-600 hover:underline">Admin &gt; Matières</NuxtLink>,
          vous pouvez créer, modifier ou supprimer une matière (semestre, abréviation, nom). Une matière déjà
          utilisée dans une campagne (un tuteur l'enseigne ou un tutee y est inscrit) ne peut pas être supprimée :
          le serveur refuse la suppression pour ne pas casser les données existantes.
        </p>
      </section>

      <section class="bg-white rounded-2xl shadow p-8 space-y-4">
        <h2 class="text-xl font-semibold text-gray-900">🤝 Affectations</h2>
        <p class="text-gray-600">
          Depuis la page « Affectations » d'une campagne, le bouton « Générer les affectations » déclenche
          l'algorithme de couplage stable (Gale-Shapley) qui associe automatiquement chaque tutee à un tuteur en
          fonction des disponibilités et des préférences de matière. Une affectation générée automatiquement ou
          créée manuellement peut ensuite être supprimée depuis cette même page si besoin.
        </p>
      </section>

      <section class="bg-white rounded-2xl shadow p-8 space-y-4">
        <h2 class="text-xl font-semibold text-gray-900">🕵️ Anonymisation des comptes</h2>
        <p class="text-gray-600">
          Un utilisateur qui quitte le cycle STPI (fin d'études, réorientation...) ne reçoit aucune notification côté plateforme : on ne peut le déduire que de son inactivité. Dans
          <NuxtLink to="/admin/users" class="text-blue-600 hover:underline">Admin &gt; Utilisateurs</NuxtLink>, la
          colonne « Statut » affiche <strong>« A quitté le cycle STPI »</strong> pour tout compte sans connexion
          depuis :
        </p>
        <ul class="list-disc list-inside text-gray-600 space-y-1">
          <li>24 mois pour un compte vu pour la dernière fois en <strong>STPI1</strong> (il lui restait potentiellement une année STPI2 à faire) ;</li>
          <li>12 mois pour un compte vu pour la dernière fois en <strong>STPI2</strong>, ou dont l'année n'était pas renseignée.</li>
        </ul>
        <p class="text-gray-600">
          Ce signalement n'est qu'une indication : c'est un admin qui décide, au cas par cas, de cliquer sur
          « Anonymiser » sur la ligne d'un utilisateur (rien n'est automatique). L'action <strong>efface
          définitivement</strong> le nom, l'email et l'identifiant CAS du compte et retire tous ses rôles,  mais
          conserve la ligne, ainsi que ses heures, séances et inscriptions déjà enregistrées, pour ne pas fausser
          les statistiques des campagnes passées. Les seuils (24 et 12 mois) sont ajustables par la personne
          qui déploie le serveur via les variables d'environnement
          <code class="bg-gray-100 px-1.5 py-0.5 rounded text-sm">ANONYMIZATION_INACTIVITY_MONTHS_STPI1</code> et
          <code class="bg-gray-100 px-1.5 py-0.5 rounded text-sm">ANONYMIZATION_INACTIVITY_MONTHS_STPI2</code>.
        </p>
      </section>

      <section class="bg-white rounded-2xl shadow p-8 space-y-4">
        <h2 class="text-xl font-semibold text-gray-900">📧 Notifications automatiques par email</h2>
        <p class="text-gray-600">
          Deux emails sont envoyés automatiquement, sans action de votre part :
        </p>
        <ul class="list-disc list-inside text-gray-600 space-y-1">
          <li>
            <strong>Mise en relation</strong> : dès qu'un tutoré est affecté à un tuteur (que ce soit via la
            génération automatique ou une affectation manuelle, sur la page « Affectations » d'une campagne), les
            deux reçoivent un email avec les coordonnées de l'autre et un lien direct vers leur espace de tutorat.
          </li>
          <li>
            <strong>Rappel de déclaration d'heures</strong> : un tuteur dont la campagne est en cours mais qui n'a
            déclaré aucune heure reçoit un rappel par email. La vérification tourne en tâche de fond toutes les 24h ;
            un même tuteur ne reçoit pas plus d'un rappel toutes les 2 semaines, et aucun rappel n'est envoyé pour
            une campagne archivée.
          </li>
        </ul>
      </section>
    </div>
  </div>
</template>

<style scoped>
</style>
